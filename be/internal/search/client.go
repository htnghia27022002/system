package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es      *elasticsearch.Client
	enabled bool
	index   string
}

type SearchRequest struct {
	Query              string
	EntityTypes        []string
	AllowedPermissions []string
	From               int
	Size               int
}

type SearchHit struct {
	EntityType string
	EntityID   string
	Title      string
	Snippet    string
	Metadata   map[string]string
	UpdatedAt  time.Time
}

type SearchResult struct {
	Hits  []SearchHit
	Total int64
}

func NewClient(url string, enabled bool) (*Client, error) {
	if !enabled {
		return &Client{enabled: false, index: IndexName}, nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{es: es, enabled: true, index: IndexName}, nil
}

func (c *Client) Enabled() bool {
	return c.enabled
}

func (c *Client) Ping(ctx context.Context) error {
	if !c.enabled {
		return fmt.Errorf("search client disabled")
	}
	res, err := c.es.Ping(c.es.Ping.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch ping: %s", res.String())
	}
	return nil
}

func (c *Client) EnsureIndex(ctx context.Context) error {
	if !c.enabled {
		return nil
	}

	res, err := c.es.Indices.Exists([]string{c.index}, c.es.Indices.Exists.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}

	body := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"entityType":      map[string]string{"type": "keyword"},
				"entityId":        map[string]string{"type": "keyword"},
				"title":           map[string]string{"type": "text"},
				"searchableText":  map[string]string{"type": "text"},
				"permissionKeys":  map[string]string{"type": "keyword"},
				"metadata":        map[string]any{"type": "object", "enabled": true},
				"updatedAt":       map[string]string{"type": "date"},
				"semanticSummary": map[string]any{"type": "text", "index": false},
			},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	createRes, err := c.es.Indices.Create(
		c.index,
		c.es.Indices.Create.WithContext(ctx),
		c.es.Indices.Create.WithBody(bytes.NewReader(payload)),
	)
	if err != nil {
		return err
	}
	defer createRes.Body.Close()
	if createRes.IsError() && createRes.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("create index: %s", createRes.String())
	}
	return nil
}

func (c *Client) Upsert(ctx context.Context, doc Document) error {
	if !c.enabled {
		return fmt.Errorf("search client disabled")
	}

	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := c.es.Index(
		c.index,
		bytes.NewReader(payload),
		c.es.Index.WithContext(ctx),
		c.es.Index.WithDocumentID(DocumentID(doc.EntityType, doc.EntityID)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("index document: %s", res.String())
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, entityType, entityID string) error {
	if !c.enabled {
		return fmt.Errorf("search client disabled")
	}

	res, err := c.es.Delete(
		c.index,
		DocumentID(entityType, entityID),
		c.es.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("delete document: %s", res.String())
	}
	return nil
}

func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
	if !c.enabled {
		return nil, fmt.Errorf("search client disabled")
	}

	filters := []any{
		map[string]any{
			"terms": map[string]any{
				"permissionKeys": req.AllowedPermissions,
			},
		},
	}
	if len(req.EntityTypes) > 0 {
		filters = append(filters, map[string]any{
			"terms": map[string]any{
				"entityType": req.EntityTypes,
			},
		})
	}

	query := map[string]any{
		"from": req.From,
		"size": req.Size,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": filters,
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query":  req.Query,
							"fields": []string{"title", "searchableText", "metadata.*"},
						},
					},
				},
			},
		},
	}

	payload, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.index),
		c.es.Search.WithBody(bytes.NewReader(payload)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("search: %s", res.String())
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source Document `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}

	out := &SearchResult{Total: parsed.Hits.Total.Value}
	for _, hit := range parsed.Hits.Hits {
		doc := hit.Source
		out.Hits = append(out.Hits, SearchHit{
			EntityType: doc.EntityType,
			EntityID:   doc.EntityID,
			Title:      doc.Title,
			Snippet:    snippet(doc.SearchableText, req.Query),
			Metadata:   doc.Metadata,
			UpdatedAt:  doc.UpdatedAt,
		})
	}
	return out, nil
}

func snippet(text, query string) string {
	text = strings.TrimSpace(text)
	if text == "" || strings.TrimSpace(query) == "" {
		return text
	}
	lower := strings.ToLower(text)
	term := strings.ToLower(strings.Fields(query)[0])
	idx := strings.Index(lower, term)
	if idx < 0 {
		if len(text) > 120 {
			return text[:120] + "..."
		}
		return text
	}
	start := idx - 20
	if start < 0 {
		start = 0
	}
	end := idx + len(term) + 80
	if end > len(text) {
		end = len(text)
	}
	return strings.TrimSpace(text[start:end])
}
