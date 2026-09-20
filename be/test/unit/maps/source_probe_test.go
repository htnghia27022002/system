package maps_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	apperrors "be/common/errors"
	mapsdto "be/internal/dto/maps"
	sourcesvc "be/internal/services/maps/source"
	"be/test/testutil"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestProbeSourceReturnsJSONBody(t *testing.T) {
	t.Parallel()
	svc := sourcesvc.NewServiceWithHTTP(
		testutil.NewMemoryDataSourceRepo(),
		roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://feeds.example/list.json" {
				t.Fatalf("url: %s", req.URL.String())
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`[{"name":"Cafe"}]`)),
				Header:     make(http.Header),
			}, nil
		}),
	)
	out, err := svc.Probe(context.Background(), mapsdto.ProbeSourceRequest{
		HTTPMethod: "GET",
		URL:        "https://feeds.example/list.json",
	})
	if err != nil {
		t.Fatalf("probe: %v", err)
	}
	if out.Status != 200 {
		t.Fatalf("status: %d", out.Status)
	}
	list, ok := out.Body.([]any)
	if !ok || len(list) != 1 {
		t.Fatalf("body: %#v", out.Body)
	}
}

func TestProbeSourceRejectsHTML(t *testing.T) {
	t.Parallel()
	svc := sourcesvc.NewServiceWithHTTP(
		testutil.NewMemoryDataSourceRepo(),
		roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<!doctype html><html></html>")),
				Header:     make(http.Header),
			}, nil
		}),
	)
	_, err := svc.Probe(context.Background(), mapsdto.ProbeSourceRequest{
		HTTPMethod: "GET",
		URL:        "https://feeds.example/page",
	})
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
}
