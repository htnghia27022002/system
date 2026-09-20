package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/common/utils"
	mapsdto "be/internal/dto/maps"
	dsmodel "be/internal/models/datasource"
	"be/internal/repository/interfaces"
)

const (
	probeTimeout = 15 * time.Second
	probeMaxBody = 5 << 20
)

// HTTPDoer fetches a probe URL (tests inject a fake).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Service implements Data Source CRUD and field-mapping validation.
type Service struct {
	sources interfaces.DataSourceRepository
	http    HTTPDoer
}

func NewService(sources interfaces.DataSourceRepository) *Service {
	return NewServiceWithHTTP(sources, nil)
}

func NewServiceWithHTTP(sources interfaces.DataSourceRepository, httpDoer HTTPDoer) *Service {
	if httpDoer == nil {
		httpDoer = &http.Client{Timeout: probeTimeout}
	}
	return &Service{sources: sources, http: httpDoer}
}

func (s *Service) List(ctx context.Context, q mapsdto.SourceListQuery) (*mapsdto.SourceListResponse, error) {
	page, limit := mapsdto.NormalizeList(q.Page, q.Limit)
	rows, total, err := s.sources.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]mapsdto.SourceRecord, 0, len(rows))
	for i := range rows {
		items = append(items, toSourceRecord(&rows[i]))
	}
	return &mapsdto.SourceListResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: int64(page*limit) < total,
	}, nil
}

func (s *Service) Get(ctx context.Context, id string) (*mapsdto.SourceRecord, error) {
	src, err := s.sources.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, apperrors.ErrNotFound
	}
	rec := toSourceRecord(src)
	return &rec, nil
}

func (s *Service) Create(ctx context.Context, userID string, req mapsdto.CreateSourceRequest) (*mapsdto.SourceRecord, error) {
	src, err := validateCreate(req)
	if err != nil {
		return nil, err
	}
	actor := utils.StringPtr(userID)
	now := time.Now().UTC()
	src.CreatedBy = actor
	src.UpdatedBy = actor
	src.CreatedAt = now
	src.UpdatedAt = now
	if err := s.sources.Create(ctx, src); err != nil {
		return nil, err
	}
	rec := toSourceRecord(src)
	return &rec, nil
}

func (s *Service) Patch(ctx context.Context, id, userID string, req mapsdto.PatchSourceRequest) (*mapsdto.SourceRecord, error) {
	src, err := s.sources.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, apperrors.ErrNotFound
	}
	if err := applyPatch(src, req); err != nil {
		return nil, err
	}
	src.UpdatedBy = utils.StringPtr(userID)
	src.UpdatedAt = time.Now().UTC()
	if err := s.sources.Update(ctx, src); err != nil {
		return nil, err
	}
	rec := toSourceRecord(src)
	return &rec, nil
}

func (s *Service) Probe(ctx context.Context, req mapsdto.ProbeSourceRequest) (*mapsdto.ProbeSourceResponse, error) {
	method, ok := dsmodel.MethodFromJSON(req.HTTPMethod)
	if !ok {
		return nil, fmt.Errorf("%w: httpMethod must be GET, POST, PUT, or PATCH", apperrors.ErrBadRequest)
	}
	if err := validateAbsoluteURL(req.URL); err != nil {
		return nil, err
	}
	var bodyReader io.Reader
	if method != dsmodel.HTTPGet && strings.TrimSpace(req.Body) != "" {
		bodyReader = strings.NewReader(req.Body)
	}
	httpReq, err := http.NewRequestWithContext(
		ctx,
		dsmodel.MethodToJSON(method),
		strings.TrimSpace(req.URL),
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", apperrors.ErrBadRequest, err.Error())
	}
	for key, value := range req.Headers {
		if s, ok := value.(string); ok && s != "" {
			httpReq.Header.Set(key, s)
		}
	}
	q := httpReq.URL.Query()
	for key, value := range req.QueryParams {
		if s, ok := value.(string); ok {
			q.Set(key, s)
		}
	}
	httpReq.URL.RawQuery = q.Encode()

	resp, err := s.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: could not reach source (%s)", apperrors.ErrBadRequest, err.Error())
	}
	defer resp.Body.Close()
	limited := io.LimitReader(resp.Body, probeMaxBody+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("%w: could not read source response", apperrors.ErrBadRequest)
	}
	if len(raw) > probeMaxBody {
		return nil, fmt.Errorf("%w: source response exceeds 5 MiB", apperrors.ErrBadRequest)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: source returned HTTP %d", apperrors.ErrBadRequest, resp.StatusCode)
	}
	trimmed := strings.TrimSpace(string(raw))
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "<!doctype html") || strings.HasPrefix(lower, "<html") {
		return nil, fmt.Errorf("%w: source response is HTML, expected JSON", apperrors.ErrBadRequest)
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("%w: source response is not JSON", apperrors.ErrBadRequest)
	}
	return &mapsdto.ProbeSourceResponse{Status: resp.StatusCode, Body: decoded}, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	src, err := s.sources.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if src == nil {
		return apperrors.ErrNotFound
	}
	return s.sources.Delete(ctx, src.ID)
}

func validateCreate(req mapsdto.CreateSourceRequest) (*dsmodel.Source, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 200 {
		return nil, fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
	}
	method, ok := dsmodel.MethodFromJSON(req.HTTPMethod)
	if !ok {
		return nil, fmt.Errorf("%w: httpMethod must be GET, POST, PUT, or PATCH", apperrors.ErrBadRequest)
	}
	if err := validateAbsoluteURL(req.URL); err != nil {
		return nil, err
	}
	if err := mapsdto.ValidateFieldMapping(req.FieldMapping); err != nil {
		return nil, fmt.Errorf("%w: %s", apperrors.ErrBadRequest, err.Error())
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return &dsmodel.Source{
		Name:         name,
		Enabled:      enabled,
		HTTPMethod:   method,
		URL:          strings.TrimSpace(req.URL),
		Headers:      stringMap(req.Headers),
		QueryParams:  stringMap(req.QueryParams),
		Body:         req.Body,
		FieldMapping: req.FieldMapping,
	}, nil
}

func applyPatch(src *dsmodel.Source, req mapsdto.PatchSourceRequest) error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len(name) > 200 {
			return fmt.Errorf("%w: name is required (max 200)", apperrors.ErrBadRequest)
		}
		src.Name = name
	}
	if req.Enabled != nil {
		src.Enabled = *req.Enabled
	}
	if req.HTTPMethod != nil {
		method, ok := dsmodel.MethodFromJSON(*req.HTTPMethod)
		if !ok {
			return fmt.Errorf("%w: httpMethod must be GET, POST, PUT, or PATCH", apperrors.ErrBadRequest)
		}
		src.HTTPMethod = method
	}
	if req.URL != nil {
		if err := validateAbsoluteURL(*req.URL); err != nil {
			return err
		}
		src.URL = strings.TrimSpace(*req.URL)
	}
	if req.Headers != nil {
		src.Headers = stringMap(req.Headers)
	}
	if req.QueryParams != nil {
		src.QueryParams = stringMap(req.QueryParams)
	}
	if req.Body != nil {
		src.Body = *req.Body
	}
	if req.FieldMapping != nil {
		if err := mapsdto.ValidateFieldMapping(req.FieldMapping); err != nil {
			return fmt.Errorf("%w: %s", apperrors.ErrBadRequest, err.Error())
		}
		src.FieldMapping = req.FieldMapping
	}
	return nil
}

func validateAbsoluteURL(raw string) error {
	raw = strings.TrimSpace(raw)
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: url must be an absolute http(s) URL", apperrors.ErrBadRequest)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("%w: url must be an absolute http(s) URL", apperrors.ErrBadRequest)
	}
	return nil
}

func stringMap(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	return in
}

func toSourceRecord(src *dsmodel.Source) mapsdto.SourceRecord {
	headers := src.Headers
	if headers == nil {
		headers = map[string]any{}
	}
	params := src.QueryParams
	if params == nil {
		params = map[string]any{}
	}
	mapping := src.FieldMapping
	if mapping == nil {
		mapping = map[string]any{}
	}
	return mapsdto.SourceRecord{
		ID:           src.ID,
		Name:         src.Name,
		Enabled:      src.Enabled,
		HTTPMethod:   dsmodel.MethodToJSON(src.HTTPMethod),
		URL:          src.URL,
		Headers:      headers,
		QueryParams:  params,
		Body:         src.Body,
		FieldMapping: mapping,
		CreatedBy:    src.CreatedBy,
		UpdatedBy:    src.UpdatedBy,
		CreatedAt:    utils.FormatRFC3339(src.CreatedAt),
		UpdatedAt:    utils.FormatRFC3339(src.UpdatedAt),
	}
}
