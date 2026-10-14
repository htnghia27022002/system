package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	authdto "be/internal/dto/auth"
	authsvc "be/internal/services/auth"
)

type AuthHandler struct {
	auth   *authsvc.Service
	oauth  *authsvc.OAuthService
	cookie RefreshCookieConfig
}

func NewAuthHandler(auth *authsvc.Service, oauth *authsvc.OAuthService, cookie RefreshCookieConfig) *AuthHandler {
	return &AuthHandler{auth: auth, oauth: oauth, cookie: cookie}
}

// issue writes a successful sign-in: in cookie mode the refresh token moves from
// the JSON body into the httpOnly cookie.
func (h *AuthHandler) issue(c *gin.Context, status int, result *authdto.AuthResponse, err error) {
	if err == nil && result != nil && h.cookie.Enabled {
		h.cookie.set(c, result.RefreshToken)
		result.RefreshToken = ""
	}
	if status == http.StatusCreated {
		httpx.Created(c, result, err)
		return
	}
	httpx.OK(c, result, err)
}

func requestClientMeta(c *gin.Context) (ip, userAgent string) {
	return httpx.ClientIP(c), c.GetHeader("User-Agent")
}

func currentSessionID(c *gin.Context) string {
	return strings.TrimSpace(c.GetHeader("X-Session-Id"))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	ip, ua := requestClientMeta(c)
	result, err := h.auth.Login(c.Request.Context(), req.Email, req.Password, ip, ua)
	h.issue(c, http.StatusOK, result, err)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	ip, ua := requestClientMeta(c)
	result, err := h.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password, ip, ua)
	h.issue(c, http.StatusCreated, result, err)
}

// refreshTokenFrom returns the refresh token from the cookie (browsers) or the JSON
// body (other clients), and whether it came from the cookie.
func (h *AuthHandler) refreshTokenFrom(c *gin.Context) (string, bool) {
	if token := h.cookie.read(c); token != "" {
		return token, true
	}
	var req authdto.RefreshRequest
	_ = c.ShouldBindJSON(&req)
	return strings.TrimSpace(req.RefreshToken), false
}

// Refresh rotates the token pair. The new refresh token goes back through the
// channel the old one came in on: cookie → cookie, JSON body → JSON body.
func (h *AuthHandler) Refresh(c *gin.Context) {
	token, fromCookie := h.refreshTokenFrom(c)
	if token == "" {
		response.Error(c, http.StatusUnauthorized, "refresh token is required")
		return
	}
	if fromCookie && c.GetHeader(CSRFHeader) == "" {
		response.Error(c, http.StatusForbidden, CSRFHeader+" header is required")
		return
	}

	ip, ua := requestClientMeta(c)
	result, err := h.auth.Refresh(c.Request.Context(), token, ip, ua)
	if fromCookie {
		if err != nil {
			h.cookie.clear(c)
		} else {
			h.cookie.set(c, result.RefreshToken)
			result.RefreshToken = ""
		}
	}
	httpx.OK(c, result, err)
}

// Logout revokes the refresh token (cookie or body) and always clears the cookie.
func (h *AuthHandler) Logout(c *gin.Context) {
	token, fromCookie := h.refreshTokenFrom(c)
	if fromCookie && c.GetHeader(CSRFHeader) == "" {
		response.Error(c, http.StatusForbidden, CSRFHeader+" header is required")
		return
	}
	if h.cookie.Enabled {
		h.cookie.clear(c)
	}

	if token != "" {
		if err := h.auth.Logout(c.Request.Context(), token); err != nil {
			response.HandleError(c, err)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.auth.Me(c.Request.Context(), userID)
	httpx.OK(c, result, err)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var req authdto.UpdateProfileRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	result, err := h.auth.UpdateProfile(c.Request.Context(), userID, req)
	httpx.OK(c, result, err)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var req authdto.ChangePasswordRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if err := h.auth.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, authdto.MessageResponse{Message: "password updated"})
}

func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	file, header, ok := httpx.FormFile(c, "avatar file is required (field: file or avatar)", "file", "avatar")
	if !ok {
		return
	}
	defer file.Close()

	result, err := h.auth.UploadAvatar(c.Request.Context(), userID, file, header)
	httpx.OK(c, result, err)
}

func (h *AuthHandler) OAuthStart(c *gin.Context) {
	provider := c.Param("provider")
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		response.Error(c, http.StatusBadRequest, "redirect_uri is required")
		return
	}
	if _, err := url.ParseRequestURI(redirectURI); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid redirect_uri")
		return
	}

	startURL, err := h.oauth.StartURL(provider, redirectURI)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.Redirect(http.StatusFound, startURL)
}

func (h *AuthHandler) OAuthProviders(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{
		"providers": h.oauth.ListConfiguredProviders(),
	})
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	var req authdto.OAuthCallbackRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	ip, ua := requestClientMeta(c)
	result, err := h.oauth.Callback(c.Request.Context(), provider, req.Code, req.RedirectURI, ip, ua)
	h.issue(c, http.StatusOK, result, err)
}

func (h *AuthHandler) ListSessions(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	result, err := h.auth.ListSessions(c.Request.Context(), userID, currentSessionID(c))
	httpx.OK(c, result, err)
}

func (h *AuthHandler) RevokeSession(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	httpx.NoContent(c, h.auth.RevokeSession(
		c.Request.Context(),
		userID,
		c.Param("id"),
		currentSessionID(c),
	))
}

func (h *AuthHandler) RevokeOtherSessions(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var req authdto.RevokeOthersRequest
	_ = c.ShouldBindJSON(&req)
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = currentSessionID(c)
	}
	result, err := h.auth.RevokeOtherSessions(c.Request.Context(), userID, sessionID)
	httpx.OK(c, result, err)
}
