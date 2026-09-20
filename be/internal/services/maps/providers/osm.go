package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"be/internal/config"
	"be/internal/services/maps"
)

const (
	defaultNominatimURL = "https://nominatim.openstreetmap.org"
	defaultUserAgent    = "system-maps/1.0 (admin search)"
	nominatimTimeout    = 8 * time.Second
	nominatimMaxBody    = 1 << 20
	nominatimLimit      = 6
)

func init() {
	maps.Register("osm", func(cfg config.Config, httpDoer maps.HTTPDoer) maps.Provider {
		return NewOSMProvider(cfg, httpDoer)
	})
}

type nominatimRow struct {
	PlaceID     any    `json:"place_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

// OSMProvider searches OpenStreetMap Nominatim.
type OSMProvider struct {
	baseURL   string
	userAgent string
	http      maps.HTTPDoer
}

func NewOSMProvider(cfg config.Config, httpDoer maps.HTTPDoer) *OSMProvider {
	if httpDoer == nil {
		httpDoer = &http.Client{Timeout: nominatimTimeout}
	}
	base := strings.TrimRight(strings.TrimSpace(cfg.MapNominatimURL), "/")
	if base == "" {
		base = defaultNominatimURL
	}
	ua := strings.TrimSpace(cfg.MapNominatimUserAgent)
	if ua == "" {
		ua = defaultUserAgent
	}
	return &OSMProvider{baseURL: base, userAgent: ua, http: httpDoer}
}

func (p *OSMProvider) ID() string { return "osm" }

func (p *OSMProvider) Search(ctx context.Context, query, countryCode string) ([]maps.Hit, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil
	}

	endpoint, err := url.Parse(p.baseURL + "/search")
	if err != nil {
		return nil, fmt.Errorf("nominatim url: %w", err)
	}
	values := endpoint.Query()
	values.Set("format", "jsonv2")
	values.Set("q", q)
	values.Set("limit", strconv.Itoa(nominatimLimit))
	if code := strings.ToLower(strings.TrimSpace(countryCode)); code != "" {
		values.Set("countrycodes", code)
	}
	endpoint.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", p.userAgent)

	res, err := p.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nominatim: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, nominatimMaxBody+1))
	if err != nil {
		return nil, fmt.Errorf("nominatim read: %w", err)
	}
	if len(body) > nominatimMaxBody {
		return nil, fmt.Errorf("nominatim: response too large")
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("nominatim: status %d", res.StatusCode)
	}

	var rows []nominatimRow
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, fmt.Errorf("nominatim json: %w", err)
	}
	return mapNominatimHits(rows), nil
}

func mapNominatimHits(rows []nominatimRow) []maps.Hit {
	hits := make([]maps.Hit, 0, len(rows))
	for _, row := range rows {
		lat, errLat := strconv.ParseFloat(strings.TrimSpace(row.Lat), 64)
		lng, errLng := strconv.ParseFloat(strings.TrimSpace(row.Lon), 64)
		if errLat != nil || errLng != nil {
			continue
		}
		display := strings.TrimSpace(row.DisplayName)
		title := strings.TrimSpace(row.Name)
		if title == "" {
			title, _, _ = strings.Cut(display, ",")
			title = strings.TrimSpace(title)
		}
		if title == "" {
			continue
		}
		subtitle := display
		if subtitle == title {
			subtitle = ""
		}
		hits = append(hits, maps.Hit{
			ExternalID: fmt.Sprint(row.PlaceID),
			Title:      title,
			Subtitle:   subtitle,
			Lat:        lat,
			Lng:        lng,
		})
	}
	return hits
}
