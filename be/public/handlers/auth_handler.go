package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	apperrors "be/internal/common/errors"
	"be/internal/common/response"
	authdto "be/internal/dto/auth"
	"be/internal/middleware"
	authsvc "be/internal/services/auth"
)

type AuthHandler struct {
	auth  *authsvc.Service
	oauth *authsvc.OAuthService
}

func NewAuthHandler(auth *authsvc.Service, oauth *authsvc.OAuthService) *AuthHandler {
	return &AuthHandler{auth: auth, oauth: oauth}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.auth.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, result)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req authdto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req authdto.LogoutRequest
	_ = c.ShouldBindJSON(&req)

	if req.RefreshToken != "" {
		if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
			response.HandleError(c, err)
			return
		}
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return
	}

	result, err := h.auth.Me(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.oauth.Callback(c.Request.Context(), provider, req.Code, req.RedirectURI)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}
