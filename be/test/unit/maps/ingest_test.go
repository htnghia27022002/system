package maps_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	apperrors "be/common/errors"
	dsmodel "be/internal/models/datasource"
	placemodel "be/internal/models/place"
	ingestsvc "be/internal/services/maps/ingest"
	"be/test/testutil"
)

type fakeHTTP struct {
	bodies map[string]string
	calls  []string
}

func (f *fakeHTTP) Do(req *http.Request) (*http.Response, error) {
	f.calls = append(f.calls, req.URL.String())
	body, ok := f.bodies[req.URL.String()]
	if !ok {
		return nil, errors.New("unexpected URL " + req.URL.String())
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}, nil
}

type noopPublisher struct{}

func (noopPublisher) Enabled() bool { return false }
func (noopPublisher) PublishMapsIngest(context.Context, string) error {
	return nil
}

type mapsStack struct {
	ingest  *ingestsvc.Service
	sources *testutil.MemoryDataSourceRepo
	runs    *testutil.MemoryIngestRunRepo
	places  *testutil.MemoryPlaceRepo
	news    *testutil.MemoryNewsRepo
	http    *fakeHTTP
}

func newMapsStack(httpFake *fakeHTTP) *mapsStack {
	if httpFake == nil {
		httpFake = &fakeHTTP{bodies: map[string]string{}}
	}
	sources := testutil.NewMemoryDataSourceRepo()
	runs := testutil.NewMemoryIngestRunRepo()
	locations := testutil.NewMemoryLocationRepo()
	places := testutil.NewMemoryPlaceRepo()
	news := testutil.NewMemoryNewsRepo()
	categories := testutil.NewMemoryCategoryRepo()
	countries := testutil.NewMemoryCountryRepo()
	divisions := testutil.NewMemoryDivisionRepo()
	svc := ingestsvc.NewService(sources, runs, locations, places, news, categories, countries, divisions, noopPublisher{}, httpFake)
	return &mapsStack{ingest: svc, sources: sources, runs: runs, places: places, news: news, http: httpFake}
}

func defaultMapping() map[string]any {
	return map[string]any{
		"location": map[string]any{
			"name":        "building.name",
			"locationKey": "building.id",
			"formatted":   "building.address",
			"lat":         "building.lat",
			"lng":         "building.lng",
		},
		"place": map[string]any{
			"name":     "shop.name",
			"category": "shop.type",
			"placeKey": "shop.id",
			"unit":     "shop.unit",
			"lat":      "shop.lat",
			"lng":      "shop.lng",
		},
		"news": map[string]any{
			"title":       "article.title",
			"originalUrl": "article.url",
		},
	}
}

func addSource(t *testing.T, stack *mapsStack, id, name, url string, enabled bool, body string) {
	t.Helper()
	src := &dsmodel.Source{
		ID:           id,
		Name:         name,
		Enabled:      enabled,
		HTTPMethod:   dsmodel.HTTPGet,
		URL:          url,
		FieldMapping: defaultMapping(),
		Headers:      map[string]any{},
		QueryParams:  map[string]any{},
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := stack.sources.Create(context.Background(), src); err != nil {
		t.Fatalf("create source: %v", err)
	}
	if enabled {
		stack.http.bodies[url] = body
	}
}

func TestDisabledSourceIsNotFetched(t *testing.T) {
	t.Parallel()
	httpFake := &fakeHTTP{bodies: map[string]string{}}
	stack := newMapsStack(httpFake)
	addSource(t, stack, "src-on", "On", "https://feeds.example/on.json", true, `[{"building":{"name":"A","id":"b1","address":"1 Main","lat":10.1,"lng":106.1},"shop":{"name":"Pho","type":"eatery","id":"s1","lat":10.1,"lng":106.1},"article":{"title":"Hi","url":"https://news.example/a"}}]`)
	addSource(t, stack, "src-off", "Off", "https://feeds.example/off.json", false, `[]`)

	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	for _, call := range httpFake.calls {
		if strings.Contains(call, "off.json") {
			t.Fatalf("disabled source was fetched: %v", httpFake.calls)
		}
	}
	if len(httpFake.calls) != 1 {
		t.Fatalf("expected one fetch, got %v", httpFake.calls)
	}
}

func TestNonListPayloadFailsThatSource(t *testing.T) {
	t.Parallel()
	httpFake := &fakeHTTP{bodies: map[string]string{}}
	stack := newMapsStack(httpFake)
	addSource(t, stack, "src-obj", "Obj", "https://feeds.example/obj.json", true, `{"not":"a list"}`)

	run, _, err := stack.ingest.Start(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if run.Status != "failed" {
		t.Fatalf("expected failed run, got %s", run.Status)
	}
	if run.Error != 1 || run.Success != 0 {
		t.Fatalf("expected one source error, got success=%d error=%d", run.Success, run.Error)
	}
}

func TestDuplicateOriginalURLDoesNotInsertSecondNews(t *testing.T) {
	t.Parallel()
	body := `[{"building":{"name":"Bldg","id":"b1","address":"1 Main","lat":10.7,"lng":106.7},"shop":{"name":"Pho","type":"eatery","id":"s1","lat":10.7,"lng":106.7},"article":{"title":"Open","url":"https://news.example/same"}}]`
	httpFake := &fakeHTTP{bodies: map[string]string{}}
	stack := newMapsStack(httpFake)
	addSource(t, stack, "src-1", "Feed", "https://feeds.example/same.json", true, body)

	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("first start: %v", err)
	}
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("second start: %v", err)
	}
	if len(stack.news.ByID) != 1 {
		t.Fatalf("expected 1 news row, got %d", len(stack.news.ByID))
	}
}

func TestUnknownCategoryDoesNotInsertSecondCategory(t *testing.T) {
	t.Parallel()
	body := `[{"building":{"name":"Bldg","id":"b2","address":"2 Main","lat":10.8,"lng":106.8},"shop":{"name":"Shop","type":"unknown-cat","id":"s2","lat":10.8,"lng":106.8},"article":{"title":"X","url":"https://news.example/x"}}]`
	stack := newMapsStack(&fakeHTTP{bodies: map[string]string{}})
	addSource(t, stack, "src-cat", "Feed", "https://feeds.example/cat.json", true, body)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if len(stack.places.ByID) != 1 {
		t.Fatalf("expected 1 place, got %d", len(stack.places.ByID))
	}
	for _, p := range stack.places.ByID {
		if p.CategoryID != "cat-uncategorized" {
			t.Fatalf("expected uncategorized, got %s", p.CategoryID)
		}
	}
}

func TestMissingCoordsCreatesPendingPlace(t *testing.T) {
	t.Parallel()
	body := `[{"building":{"name":"Bldg","id":"b3","address":"3 Main"},"shop":{"name":"NoPin","type":"hotel","id":"s3"},"article":{"title":"Y","url":"https://news.example/y"}}]`
	stack := newMapsStack(&fakeHTTP{bodies: map[string]string{}})
	addSource(t, stack, "src-pend", "Feed", "https://feeds.example/pend.json", true, body)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if len(stack.places.ByID) != 1 {
		t.Fatalf("expected 1 place, got %d", len(stack.places.ByID))
	}
	for _, p := range stack.places.ByID {
		if p.Status != placemodel.StatusPending {
			t.Fatalf("expected pending, got %d", p.Status)
		}
	}
}

func TestConcurrentStartReturnsConflict(t *testing.T) {
	t.Parallel()
	stack := newMapsStack(&fakeHTTP{bodies: map[string]string{}})
	_ = stack.runs.Insert(context.Background(), &dsmodel.IngestRun{
		ID:          "run-active",
		Status:      dsmodel.RunRunning,
		TriggeredBy: "user-1",
		CreatedAt:   time.Now().UTC(),
	})
	_, _, err := stack.ingest.Start(context.Background(), "user-1")
	if err == nil || !apperrors.IsConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestUnsavedURLIsNeverFetched(t *testing.T) {
	t.Parallel()
	httpFake := &fakeHTTP{bodies: map[string]string{
		"https://feeds.example/saved.json": `[]`,
	}}
	stack := newMapsStack(httpFake)
	addSource(t, stack, "src-saved", "Saved", "https://feeds.example/saved.json", true, `[]`)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	for _, call := range httpFake.calls {
		if call == "https://open-web.example/not-saved.json" {
			t.Fatal("fetched a URL that is not a saved source")
		}
	}
	if len(httpFake.calls) != 1 || httpFake.calls[0] != "https://feeds.example/saved.json" {
		t.Fatalf("unexpected fetches: %v", httpFake.calls)
	}
}
