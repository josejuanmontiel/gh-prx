package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type AgilityIssueQuery struct {
	From   string                 `json:"from"`
	Select []string               `json:"select"`
	Where  map[string]interface{} `json:"where",omitempty`
	Sort   []string               `json:"sort"`
}

type AgilityIssues [][]AgilityIssue

type Oid struct {
	Oid string `json:"_oid"`
}

type AgilityIssue struct {
	Oid          string `json:"_oid"`
	ID           Oid    `json:"ID"`
	Number       string `json:"Number"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	Timebox      Oid    `json:"Timebox"`
	Parent       Oid    `json:"Parent"`
	ParentNumber string `json:"Parent.Number"`
	ParentName   string `json:"Parent.Name"`
}

func (p *AgilityIssueProvider) Get(ctx context.Context, id string) (*models.Issue, error) {

	query := AgilityIssueQuery{
		From: "Task",
		Select: []string{
			"ID",
			"Number",
			"Name",
			"Description",
			"Owners",
			"Timebox",
			"Parent",
			"Parent.Number",
			"Parent.Name",
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
	query := AgilityIssueQuery{
		From: "Task",
		Select: []string{
			"ID",
			"Number",
			"Name",
			"Status.Name",
			"Description",
			"Owners",
			"Timebox",
			"Parent",
			"Parent.Number",
			"Parent.Name",
			"CreateDate",
		},
		Where: map[string]interface{}{},
		Sort: []string{
			"-CreateDate",
		},
	}

	if p.Config.Owner != "" {
		query.Where = map[string]interface{}{
			"Owners":      p.Config.Owner,
			"Status.Name": "In Progress",
		}
	}

	issue := AgilityIssues{}
	if err := p.query(ctx, "/query.v1", query, &issue); err != nil {
		return nil, err
	}

	result := make([]*models.Issue, len(issue[0]))
	for i, issue := range issue[0] {
		result[i] = issue.ToIssue()
	}

	return result, nil
}

func (i *AgilityIssue) ToIssue() *models.Issue {
	issueType := LabelToType["enhancement"]

	return &models.Issue{
		Key:    i.Number,
		Title:  i.Name,
		Parent: i.ParentNumber,
		Type:   issueType,
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
