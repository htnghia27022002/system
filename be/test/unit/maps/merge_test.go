package maps_test

import (
	"context"
	"testing"

	placemodel "be/internal/models/place"
)

func TestTwoShopsSameCoordsStayTwoPlaces(t *testing.T) {
	t.Parallel()
	body := `[
		{"building":{"name":"Chung cu","id":"bldg-1","address":"12 Nguyen Hue","lat":10.7769,"lng":106.7009},"shop":{"name":"Pho A","type":"eatery","id":"shop-a","lat":10.7769,"lng":106.7009},"article":{"title":"A","url":"https://news.example/a"}},
		{"building":{"name":"Chung cu","id":"bldg-1","address":"12 Nguyen Hue","lat":10.7769,"lng":106.7009},"shop":{"name":"Cafe B","type":"eatery","id":"shop-b","lat":10.7769,"lng":106.7009},"article":{"title":"B","url":"https://news.example/b"}}
	]`
	stack := newMapsStack(nil)
	addSource(t, stack, "src-two", "Feed", "https://feeds.example/two.json", true, body)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if len(stack.places.ByID) != 2 {
		t.Fatalf("expected 2 places, got %d", len(stack.places.ByID))
	}
	if len(stack.news.ByID) != 2 {
		t.Fatalf("expected 2 news rows, got %d", len(stack.news.ByID))
	}
}

func TestSamePlaceKeyMergesAndAddsNews(t *testing.T) {
	t.Parallel()
	first := `[{"building":{"name":"Chung cu","id":"bldg-2","address":"1 Street","lat":10.1,"lng":106.1},"shop":{"name":"Shop","type":"restaurant","id":"shop-1","lat":10.1,"lng":106.1},"article":{"title":"One","url":"https://news.example/one"}}]`
	second := `[{"building":{"name":"Chung cu","id":"bldg-2","address":"1 Street","lat":10.1,"lng":106.1},"shop":{"name":"Shop","type":"restaurant","id":"shop-1","lat":10.1,"lng":106.1},"article":{"title":"Two","url":"https://news.example/two"}}]`
	stack := newMapsStack(nil)
	addSource(t, stack, "src-m1", "Feed1", "https://feeds.example/m1.json", true, first)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("first: %v", err)
	}
	addSource(t, stack, "src-m2", "Feed2", "https://feeds.example/m2.json", true, second)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("second: %v", err)
	}
	if len(stack.places.ByID) != 1 {
		t.Fatalf("expected 1 place, got %d", len(stack.places.ByID))
	}
	if len(stack.news.ByID) != 2 {
		t.Fatalf("expected 2 news rows, got %d", len(stack.news.ByID))
	}
}

func TestSameArticleURLRefreshesInPlace(t *testing.T) {
	t.Parallel()
	first := `[{"building":{"name":"Site","id":"bldg-3","address":"2 Street","lat":11.1,"lng":107.1},"shop":{"name":"Shop","type":"hotel","id":"shop-9","lat":11.1,"lng":107.1},"article":{"title":"Old","url":"https://news.example/same-url"}}]`
	second := `[{"building":{"name":"Site","id":"bldg-3","address":"2 Street","lat":11.1,"lng":107.1},"shop":{"name":"Shop","type":"hotel","id":"shop-9","lat":11.1,"lng":107.1},"article":{"title":"New","url":"https://news.example/same-url"}}]`
	stack := newMapsStack(nil)
	addSource(t, stack, "src-u1", "Feed1", "https://feeds.example/u1.json", true, first)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("first: %v", err)
	}
	addSource(t, stack, "src-u2", "Feed2", "https://feeds.example/u2.json", true, second)
	if _, _, err := stack.ingest.Start(context.Background(), "user-1"); err != nil {
		t.Fatalf("second: %v", err)
	}
	if len(stack.news.ByID) != 1 {
		t.Fatalf("expected 1 news row, got %d", len(stack.news.ByID))
	}
	for _, item := range stack.news.ByID {
		if item.Title != "New" {
			t.Fatalf("expected refreshed title, got %q", item.Title)
		}
	}
	for _, p := range stack.places.ByID {
		if p.Status != placemodel.StatusActive {
			t.Fatalf("expected active place, got %d", p.Status)
		}
	}
}
