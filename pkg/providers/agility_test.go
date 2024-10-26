package providers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilaif/gh-prx/pkg/config"
	"github.com/ilaif/gh-prx/pkg/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestAgilityIssueProvider_Get(t *testing.T) {
	// Arrange
	config := &config.AgilityConfig{
		APIKey: "test-api-key",
	}

	provider := &providers.AgilityIssueProvider{
		Config: config,
	}

	// Mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer test-api-key")
		response := [][]providers.AgilityIssue{
			{
				{
					Oid:         "123",
					Name:        "Test Issue",
					Number:      "1",
					ID:          providers.StoryID{Oid: "123"},
					Description: "This is a test issue",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	provider.Config.Endpoint = server.URL // set the mock server URL

	// Act
	issue, err := provider.Get(context.Background(), "1")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, "123", issue.Key)
	assert.Equal(t, "Test Issue", issue.Title)
}

func TestAgilityIssueProvider_List(t *testing.T) {
	// Arrange
	config := &config.AgilityConfig{
		APIKey: "test-api-key",
	}

	provider := &providers.AgilityIssueProvider{
		Config: config,
	}

	// Mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer test-api-key")

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"AgilityIssue": []map[string]interface{}{
					{
						"_oid":        "123",
						"Name":        "Test Issue 1",
						"Number":      "1",
						"ID":          map[string]interface{}{"_oid": "123"},
						"Description": "This is the first test issue",
					},
					{
						"_oid":        "456",
						"Name":        "Test Issue 2",
						"Number":      "2",
						"ID":          map[string]interface{}{"_oid": "456"},
						"Description": "This is the second test issue",
					},
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	provider.Config.APIKey = server.URL // set the mock server URL

	// Act
	issues, err := provider.List(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "Test Issue 1", issues[0].Title)
	assert.Equal(t, "Test Issue 2", issues[1].Title)
}
