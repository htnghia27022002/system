package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/common/utils"
	mapsdto "be/internal/dto/maps"
	dsmodel "be/internal/models/datasource"
	"be/internal/repository/interfaces"
)

const (
	FetchTimeout   = 30 * time.Second
	MaxBodyBytes   = 5 << 20
	noEnabledMsg   = "no enabled sources"
	inProgressMsg  = "ingest already in progress"
)

var (
	ErrIngestInProgress = fmt.Errorf("%w: %s", apperrors.ErrConflict, inProgressMsg)
	errSkipItem         = errors.New("skip item")
)

// Publisher publishes a queued ingest run when NATS is available.
type Publisher interface {
	Enabled() bool
	PublishMapsIngest(ctx context.Context, runID string) error
}

// HTTPDoer fetches a saved Source URL (tests inject a fake).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Service starts and processes Load data runs.
type Service struct {
	sources   interfaces.DataSourceRepository
	runs      interfaces.IngestRunRepository
	locations interfaces.LocationRepository
	places    interfaces.PlaceRepository
	news      interfaces.NewsRepository
	categories interfaces.CategoryRepository
	countries interfaces.CountryRepository
	divisions interfaces.DivisionRepository
	publisher Publisher
	http      HTTPDoer
}

func NewService(
	sources interfaces.DataSourceRepository,
	runs interfaces.IngestRunRepository,
	locations interfaces.LocationRepository,
	places interfaces.PlaceRepository,
	news interfaces.NewsRepository,
	categories interfaces.CategoryRepository,
	countries interfaces.CountryRepository,
	divisions interfaces.DivisionRepository,
	publisher Publisher,
	httpDoer HTTPDoer,
) *Service {
	if httpDoer == nil {
		httpDoer = &http.Client{Timeout: FetchTimeout}
	}
	return &Service{
		sources:    sources,
		runs:       runs,
		locations:  locations,
		places:     places,
		news:       news,
		categories: categories,
		countries:  countries,
		divisions:  divisions,
		publisher:  publisher,
		http:       httpDoer,
	}
}

// Start inserts a queued run. Returns the run and whether it was processed synchronously.
func (s *Service) Start(ctx context.Context, userID string) (*mapsdto.IngestRun, bool, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, false, apperrors.ErrUnauthorized
	}
	active, err := s.runs.HasActive(ctx)
	if err != nil {
		return nil, false, err
	}
	if active {
		return nil, false, ErrIngestInProgress
	}

	enabled, err := s.sources.ListEnabled(ctx)
	if err != nil {
		return nil, false, err
	}
	actor := utils.StringPtr(userID)
	now := time.Now().UTC()
	run := &dsmodel.IngestRun{
		Status:       dsmodel.RunQueued,
		TriggeredBy:  userID,
		SourcesTotal: len(enabled),
		CreatedBy:    actor,
		UpdatedBy:    actor,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.runs.Insert(ctx, run); err != nil {
		return nil, false, err
	}

	sync := s.publisher == nil || !s.publisher.Enabled()
	if sync {
		if err := s.ProcessRun(ctx, run.ID); err != nil {
			return nil, false, err
		}
		dto, err := s.Get(ctx, run.ID)
		return dto, true, err
	}
	if err := s.publisher.PublishMapsIngest(ctx, run.ID); err != nil {
		if procErr := s.ProcessRun(ctx, run.ID); procErr != nil {
			return nil, false, procErr
		}
		dto, getErr := s.Get(ctx, run.ID)
		return dto, true, getErr
	}
	dto, err := s.Get(ctx, run.ID)
	return dto, false, err
}

func (s *Service) Latest(ctx context.Context) (*mapsdto.LatestIngestResponse, error) {
	run, err := s.runs.Latest(ctx)
	if err != nil {
		return nil, err
	}
	if run == nil {
		return &mapsdto.LatestIngestResponse{Run: nil}, nil
	}
	dto, err := s.toRunDTO(ctx, run)
	if err != nil {
		return nil, err
	}
	return &mapsdto.LatestIngestResponse{Run: dto}, nil
}

func (s *Service) Get(ctx context.Context, id string) (*mapsdto.IngestRun, error) {
	run, err := s.runs.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, apperrors.ErrNotFound
	}
	return s.toRunDTO(ctx, run)
}

// ProcessQueued drains queued runs (NATS-disabled worker poll).
func (s *Service) ProcessQueued(ctx context.Context, limit int) (int, error) {
	runs, err := s.runs.ListQueued(ctx, limit)
	if err != nil {
		return 0, err
	}
	n := 0
	for i := range runs {
		if err := s.ProcessRun(ctx, runs[i].ID); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// ProcessRun claims a queued run and fetches only enabled saved Sources.
func (s *Service) ProcessRun(ctx context.Context, runID string) error {
	run, err := s.runs.ClaimQueued(ctx, runID)
	if err != nil {
		return err
	}
	if run == nil {
		existing, getErr := s.runs.GetByID(ctx, runID)
		if getErr != nil {
			return getErr
		}
		if existing == nil {
			return apperrors.ErrNotFound
		}
		return nil
	}

	actor := run.CreatedBy
	if actor == nil {
		actor = utils.StringPtr(run.TriggeredBy)
	}
	enabled, err := s.sources.ListEnabled(ctx)
	if err != nil {
		return s.failRun(ctx, run, actor, err.Error())
	}
	run.SourcesTotal = len(enabled)
	if len(enabled) == 0 {
		return s.failRun(ctx, run, actor, noEnabledMsg)
	}

	success, failed := 0, 0
	itemsOK, itemsErr := 0, 0
	for i := range enabled {
		src := enabled[i]
		okCount, errCount, srcErr := s.processSource(ctx, run, actor, &src)
		itemsOK += okCount
		itemsErr += errCount
		if srcErr != nil {
			failed++
		} else {
			success++
		}
	}

	now := time.Now().UTC()
	run.Success = success
	run.Error = failed
	run.ItemsSuccess = itemsOK
	run.ItemsError = itemsErr
	run.FinishedAt = &now
	run.UpdatedBy = actor
	run.UpdatedAt = now
	if failed == len(enabled) {
		run.Status = dsmodel.RunFailed
		if run.ErrorMessage == "" {
			run.ErrorMessage = "all enabled sources failed"
		}
	} else {
		run.Status = dsmodel.RunCompleted
	}
	return s.runs.Update(ctx, run)
}

func (s *Service) processSource(ctx context.Context, run *dsmodel.IngestRun, actor *string, src *dsmodel.Source) (int, int, error) {
	now := time.Now().UTC()
	row := &dsmodel.IngestRunSource{
		RunID:      run.ID,
		SourceID:   utils.StringPtr(src.ID),
		SourceName: src.Name,
		Status:     dsmodel.SourceRunning,
		CreatedBy:  actor,
		UpdatedBy:  actor,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.runs.InsertSource(ctx, row); err != nil {
		return 0, 0, err
	}

	items, fetchErr := s.fetchList(ctx, src)
	if fetchErr != nil {
		row.Status = dsmodel.SourceFailed
		row.ErrorMessage = fetchErr.Error()
		row.UpdatedAt = time.Now().UTC()
		_ = s.runs.UpdateSource(ctx, row)
		return 0, 0, fetchErr
	}

	okCount, errCount := 0, 0
	for _, raw := range items {
		mapped, mapErr := mapItem(raw, src.FieldMapping)
		if mapErr != nil {
			errCount++
			continue
		}
		if err := s.persistItem(ctx, actor, utils.StringPtr(src.ID), src.Name, mapped); err != nil {
			if errors.Is(err, errSkipItem) {
				errCount++
				continue
			}
			errCount++
			continue
		}
		okCount++
	}

	row.ItemsSuccess = okCount
	row.ItemsError = errCount
	row.UpdatedBy = actor
	row.UpdatedAt = time.Now().UTC()
	if okCount == 0 && errCount > 0 && len(items) > 0 {
		row.Status = dsmodel.SourceFailed
		row.ErrorMessage = "all items failed"
		_ = s.runs.UpdateSource(ctx, row)
		return okCount, errCount, errors.New(row.ErrorMessage)
	}
	row.Status = dsmodel.SourceCompleted
	_ = s.runs.UpdateSource(ctx, row)
	return okCount, errCount, nil
}

func (s *Service) fetchList(ctx context.Context, src *dsmodel.Source) ([]any, error) {
	if src == nil || !src.Enabled {
		return nil, fmt.Errorf("source is disabled")
	}
	req, err := http.NewRequestWithContext(ctx, dsmodel.MethodToJSON(src.HTTPMethod), src.URL, requestBody(src))
	if err != nil {
		return nil, err
	}
	for key, value := range src.Headers {
		if s, ok := value.(string); ok && s != "" {
			req.Header.Set(key, s)
		}
	}
	q := req.URL.Query()
	for key, value := range src.QueryParams {
		if s, ok := value.(string); ok {
			q.Set(key, s)
		}
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("source returned HTTP %d", resp.StatusCode)
	}
	limited := io.LimitReader(resp.Body, MaxBodyBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(body) > MaxBodyBytes {
		return nil, fmt.Errorf("source response exceeds 5 MiB")
	}
	if looksLikeHTML(body) {
		return nil, fmt.Errorf("source response is HTML, expected a JSON list")
	}
	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("source response is not JSON")
	}
	return extractItems(decoded, src.FieldMapping)
}

func requestBody(src *dsmodel.Source) io.Reader {
	if src.HTTPMethod == dsmodel.HTTPGet || strings.TrimSpace(src.Body) == "" {
		return nil
	}
	return strings.NewReader(src.Body)
}

func looksLikeHTML(body []byte) bool {
	trimmed := strings.TrimSpace(string(body))
	lower := strings.ToLower(trimmed)
	return strings.HasPrefix(lower, "<!doctype html") || strings.HasPrefix(lower, "<html")
}

func (s *Service) failRun(ctx context.Context, run *dsmodel.IngestRun, actor *string, message string) error {
	now := time.Now().UTC()
	run.Status = dsmodel.RunFailed
	run.ErrorMessage = message
	run.FinishedAt = &now
	run.UpdatedBy = actor
	run.UpdatedAt = now
	return s.runs.Update(ctx, run)
}

func (s *Service) toRunDTO(ctx context.Context, run *dsmodel.IngestRun) (*mapsdto.IngestRun, error) {
	rows, err := s.runs.ListSources(ctx, run.ID)
	if err != nil {
		return nil, err
	}
	sources := make([]mapsdto.IngestSourceOutcome, 0, len(rows))
	for i := range rows {
		row := rows[i]
		sources = append(sources, mapsdto.IngestSourceOutcome{
			ID:           row.ID,
			SourceID:     row.SourceID,
			SourceName:   row.SourceName,
			Status:       dsmodel.SourceStatusToJSON(row.Status),
			ErrorMessage: row.ErrorMessage,
			ItemsSuccess: row.ItemsSuccess,
			ItemsError:   row.ItemsError,
		})
	}
	return &mapsdto.IngestRun{
		ID:           run.ID,
		Status:       dsmodel.RunStatusToJSON(run.Status),
		ErrorMessage: run.ErrorMessage,
		SourcesTotal: run.SourcesTotal,
		Success:      run.Success,
		Error:        run.Error,
		ItemsSuccess: run.ItemsSuccess,
		ItemsError:   run.ItemsError,
		Sources:      sources,
		StartedAt:    utils.FormatRFC3339Ptr(run.StartedAt),
		FinishedAt:   utils.FormatRFC3339Ptr(run.FinishedAt),
		CreatedAt:    utils.FormatRFC3339(run.CreatedAt),
	}, nil
}
