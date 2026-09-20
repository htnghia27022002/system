package datasource

import "time"

const (
	RunQueued    int16 = 1
	RunRunning   int16 = 2
	RunCompleted int16 = 3
	RunFailed    int16 = 4
)

const (
	SourcePending   int16 = 1
	SourceRunning   int16 = 2
	SourceCompleted int16 = 3
	SourceFailed    int16 = 4
	SourceSkipped   int16 = 5
)

// IngestRun is one operator-triggered Load data execution.
type IngestRun struct {
	ID           string     `json:"id" db:"id"`
	Status       int16      `json:"status" db:"status"`
	TriggeredBy  string     `json:"triggeredBy" db:"triggered_by"`
	ErrorMessage string     `json:"errorMessage" db:"error_message"`
	SourcesTotal int        `json:"sourcesTotal" db:"sources_total"`
	Success      int        `json:"success" db:"success"`
	Error        int        `json:"error" db:"error"`
	ItemsSuccess int        `json:"itemsSuccess" db:"items_success"`
	ItemsError   int        `json:"itemsError" db:"items_error"`
	StartedAt    *time.Time `json:"startedAt" db:"started_at"`
	FinishedAt   *time.Time `json:"finishedAt" db:"finished_at"`
	CreatedBy    *string    `json:"createdBy" db:"created_by"`
	UpdatedBy    *string    `json:"updatedBy" db:"updated_by"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

func (IngestRun) TableName() string {
	return "data_ingest_runs"
}

func RunStatusToJSON(status int16) string {
	switch status {
	case RunRunning:
		return "running"
	case RunCompleted:
		return "completed"
	case RunFailed:
		return "failed"
	default:
		return "queued"
	}
}

// IngestRunSource is the per–Data Source outcome for one ingest run.
type IngestRunSource struct {
	ID           string    `json:"id" db:"id"`
	RunID        string    `json:"runId" db:"run_id"`
	SourceID     *string   `json:"sourceId" db:"source_id"`
	SourceName   string    `json:"sourceName" db:"source_name"`
	Status       int16     `json:"status" db:"status"`
	ErrorMessage string    `json:"errorMessage" db:"error_message"`
	ItemsSuccess int       `json:"itemsSuccess" db:"items_success"`
	ItemsError   int       `json:"itemsError" db:"items_error"`
	CreatedBy    *string   `json:"createdBy" db:"created_by"`
	UpdatedBy    *string   `json:"updatedBy" db:"updated_by"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

func (IngestRunSource) TableName() string {
	return "data_ingest_run_sources"
}

func SourceStatusToJSON(status int16) string {
	switch status {
	case SourceRunning:
		return "running"
	case SourceCompleted:
		return "completed"
	case SourceFailed:
		return "failed"
	case SourceSkipped:
		return "skipped"
	default:
		return "pending"
	}
}
