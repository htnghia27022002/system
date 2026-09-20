package maps

// IngestSourceOutcome is one per–Data Source result on an ingest run.
type IngestSourceOutcome struct {
	ID           string  `json:"id"`
	SourceID     *string `json:"sourceId"`
	SourceName   string  `json:"sourceName"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"errorMessage"`
	ItemsSuccess int     `json:"itemsSuccess"`
	ItemsError   int     `json:"itemsError"`
}

// IngestRun is GET/POST ingest payload including per-source outcomes.
type IngestRun struct {
	ID           string                `json:"id"`
	Status       string                `json:"status"`
	ErrorMessage string                `json:"errorMessage"`
	SourcesTotal int                   `json:"sourcesTotal"`
	Success      int                   `json:"success"`
	Error        int                   `json:"error"`
	ItemsSuccess int                   `json:"itemsSuccess"`
	ItemsError   int                   `json:"itemsError"`
	Sources      []IngestSourceOutcome `json:"sources"`
	StartedAt    *string               `json:"startedAt"`
	FinishedAt   *string               `json:"finishedAt"`
	CreatedAt    string                `json:"createdAt"`
}

// LatestIngestResponse wraps the latest run or null.
type LatestIngestResponse struct {
	Run *IngestRun `json:"run"`
}
