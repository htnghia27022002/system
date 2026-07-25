package webhook

// InboxResponse is returned by get-or-create and regenerate.
type InboxResponse struct {
	ID               string `json:"id"`
	PublicUUID       string `json:"publicUuid"`
	PublicPath       string `json:"publicPath"`
	ActiveCount      int    `json:"activeCount"`
	LifetimeReceived  int    `json:"lifetimeReceived"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// ListRequestsQuery is the owner list filter/pagination input.
type ListRequestsQuery struct {
	Method string `form:"method"`
	Q      string `form:"q"`
	// Read filter: empty/"all" = no filter; "read" | "unread".
	Read  string `form:"read"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}

// RequestListItem is one row in the owner inbox list.
type RequestListItem struct {
	ID        string `json:"id"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	ClientIP  string `json:"clientIp"`
	CreatedAt string `json:"createdAt"`
	Snippet   string `json:"snippet"`
	IsRead    bool   `json:"isRead"`
}

// RequestListResponse wraps paginated active requests plus inbox counters.
type RequestListResponse struct {
	Items            []RequestListItem `json:"items"`
	ActiveCount      int               `json:"activeCount"`
	LifetimeReceived  int               `json:"lifetimeReceived"`
	Page             int               `json:"page"`
	Limit            int               `json:"limit"`
	Total            int64             `json:"total"`
	HasMore          bool              `json:"hasMore"`
}

// RequestDetailResponse is the full owner detail payload.
type RequestDetailResponse struct {
	ID            string         `json:"id"`
	InboxID       string         `json:"inboxId"`
	Method        string         `json:"method"`
	URL           string         `json:"url"`
	ClientIP      string         `json:"clientIp"`
	Headers       map[string]any `json:"headers"`
	Query         map[string]any `json:"query"`
	Form          map[string]any `json:"form"`
	Body          string         `json:"body"`
	BodyEncoding  string         `json:"bodyEncoding"`
	IsBinary      bool           `json:"isBinary"`
	ContentType   string         `json:"contentType"`
	BodyTruncated bool           `json:"bodyTruncated"`
	CaptureStatus string         `json:"captureStatus"`
	IsRead        bool           `json:"isRead"`
	CreatedAt     string         `json:"createdAt"`
}

// SetReadRequest marks a request read or unread.
type SetReadRequest struct {
	IsRead bool `json:"isRead"`
}

// SoftDeleteResponse is returned after soft-deleting a request.
type SoftDeleteResponse struct {
	OK              bool `json:"ok"`
	ActiveCount     int  `json:"activeCount"`
	LifetimeReceived int  `json:"lifetimeReceived"`
}

// CaptureAckResponse is the fixed public capture acknowledgment.
type CaptureAckResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}
