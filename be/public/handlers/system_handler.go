package handlers

import (
	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	admindto "be/internal/dto/admin"
	"be/internal/middleware"
	auditlogsvc "be/internal/services/auditlog"
	dashboardsvc "be/internal/services/dashboard"
	"be/internal/settings"
)

// SystemHandler serves the admin audit log, settings and dashboard.
type SystemHandler struct {
	audit     *auditlogsvc.Service
	settings  *settings.Service
	dashboard *dashboardsvc.Service
}

func NewSystemHandler(audit *auditlogsvc.Service, settingsSvc *settings.Service, dashboard *dashboardsvc.Service) *SystemHandler {
	return &SystemHandler{audit: audit, settings: settingsSvc, dashboard: dashboard}
}

func (h *SystemHandler) ListAuditLogs(c *gin.Context) {
	var form admindto.ListAuditLogsQuery
	if !httpx.BindQuery(c, &form) {
		return
	}
	page, err := h.audit.List(c.Request.Context(), form)
	httpx.OK(c, page, err)
}

func (h *SystemHandler) ListSettings(c *gin.Context) {
	values, err := h.settings.List(c.Request.Context())
	httpx.OK(c, gin.H{"items": values}, err)
}

func (h *SystemHandler) UpdateSettings(c *gin.Context) {
	var req admindto.UpdateSettingsRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if _, err := h.settings.Update(c.Request.Context(), req.Values, middleware.GetUserID(c)); err != nil {
		response.HandleError(c, err)
		return
	}
	h.ListSettings(c)
}

func (h *SystemHandler) DashboardOverview(c *gin.Context) {
	overview, err := h.dashboard.Overview(c.Request.Context())
	httpx.OK(c, overview, err)
}
