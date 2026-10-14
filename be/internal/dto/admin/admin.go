package admin

import "be/pkg/query"

// ListAuditLogsQuery filters GET /api/admin/audit-logs.
type ListAuditLogsQuery struct {
	query.PageParams
	// Search matches action, actor email or resource id (case-insensitive).
	Search       string `form:"search" binding:"omitempty,max=200"`
	Action       string `form:"action" binding:"omitempty,max=100"`
	ActorID      string `form:"actorId" binding:"omitempty,uuid"`
	ResourceType string `form:"resourceType" binding:"omitempty,max=50"`
	ResourceID   string `form:"resourceId" binding:"omitempty,max=100"`
	// From / To bound created_at (RFC 3339).
	From string `form:"from" binding:"omitempty"`
	To   string `form:"to" binding:"omitempty"`
}

type AuditLogResponse struct {
	ID           string         `json:"id"`
	ActorID      *string        `json:"actorId"`
	ActorEmail   string         `json:"actorEmail"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resourceType"`
	ResourceID   string         `json:"resourceId"`
	Metadata     map[string]any `json:"metadata"`
	IPAddress    string         `json:"ipAddress"`
	UserAgent    string         `json:"userAgent"`
	RequestID    string         `json:"requestId"`
	CreatedAt    string         `json:"createdAt"`
}

// UpdateSettingsRequest is a partial update: key → new value.
type UpdateSettingsRequest struct {
	Values map[string]any `json:"values" binding:"required"`
}

type DashboardKPI struct {
	Key           string  `json:"key"`
	Value         int64   `json:"value"`
	ChangePercent float64 `json:"changePercent"`
}

type DashboardPoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type DashboardOverview struct {
	KPIs []DashboardKPI `json:"kpis"`
	// WeeklyActivity is successful sign-ins per UTC day, oldest first (7 days).
	WeeklyActivity []DashboardPoint `json:"weeklyActivity"`
	UpdatedAt      string           `json:"updatedAt"`
}
