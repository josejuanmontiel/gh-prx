package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ilaif/gh-prx/pkg/config"
	"github.com/ilaif/gh-prx/pkg/models"
)

type AgilityIssueProvider struct {
	Config *config.AgilityConfig
}

func (p *AgilityIssueProvider) Name() string {
	return "agility"
}

func (p *AgilityIssueProvider) Get(ctx context.Context, id string) (*models.Issue, error) {

	query := AgilityIssueQuery{
		From: "Story",
		Select: []string{
			"Name",
			"Number",
			"ID",
			"Description",
		},
		Where: map[string]interface{}{
			"Number": id,
		},
	}

	issue := AgilityIssues{}
	if err := p.query(ctx, "/query.v1", query, &issue); err != nil {
		return nil, err
	}

	return issue[0][0].ToIssue(), nil
}

func (p *AgilityIssueProvider) List(ctx context.Context) ([]*models.Issue, error) {
	query := AgilityIssues{}
	// WIP
	// issues := &AgilityIssues{}
	// if err := p.query(ctx, path, issues); err != nil {
	// 	return nil, err
	// }

	result := make([]*models.Issue, len(query[0]))
	for i, issue := range query[0] {
		result[i] = issue.ToIssue()
	}

	return result, nil
}

type AgilityIssueQuery struct {
	From   string                 `json:"from"`
	Select []string               `json:"select"`
	Where  map[string]interface{} `json:"where"`
}

type AgilityIssues [][]AgilityIssue

type StoryID struct {
	Oid string `json:"_oid"`
}

type AgilityIssue struct {
	Oid         string  `json:"_oid"`
	Name        string  `json:"Name"`
	Number      string  `json:"Number"`
	ID          StoryID `json:"ID"`
	Description string  `json:"Description"`
}

func (i *AgilityIssue) ToIssue() *models.Issue {
	// issueType := ""
	// for _, label := range i.Labels.Nodes {
	// 	if it, ok := LabelToType[strings.ToLower(label.Name)]; ok {
	// 		issueType = it

	// 		break
	// 	}
	// }

	// Expresión regular para eliminar etiquetas HTML
	re := regexp.MustCompile(`<[^>]*>`)
	plain := re.ReplaceAllString(i.Description, "")
	plain = strings.ReplaceAll(plain, " ", "_")

	return &models.Issue{
		Key:   i.ID.Oid,
		Title: i.Name,
		// Type:                issueType,
		SuggestedBranchName: plain,
	}
}

func (p *AgilityIssueProvider) query(ctx context.Context, path string, body any, response any) error {
	url := fmt.Sprintf("%s/%s", p.Config.Endpoint, path)

	// Serialización de la estructura a JSON
	payload, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse body '%s'", url)
	}

	// Imprimir el body en formato JSON
	fmt.Printf("Request Body: %s\n", string(payload)) // Añadir esta línea

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrapf(err, "Failed to create request for '%s'", url)
	}
	req.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	req.Header.Add("content-type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to request for '%s'", url)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return errors.Errorf("Request '%s' not found", path)
		}

		return errors.Errorf("Request '%s' failed: %s", path, res.Status)
	}

	// Leer el cuerpo de la respuesta antes de decodificarlo
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body")
	}

	// Imprimir el cuerpo de la respuesta
	fmt.Printf("Response Body: %s\n", responseBody) // Añadir esta línea

	// Decodificar el cuerpo de la respuesta en la estructura de respuesta
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return errors.Wrap(err, "Failed to parse response")
	}

	return nil
}
