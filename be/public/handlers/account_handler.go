package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	accountdto "be/internal/dto/account"
	"be/internal/middleware"
	accountsvc "be/internal/services/account"
	usersvc "be/internal/services/user"
)

// AccountHandler serves the email-driven account flows (reset, verify, invite).
type AccountHandler struct {
	svc   *accountsvc.Service
	users *usersvc.Service
}

func NewAccountHandler(svc *accountsvc.Service, users *usersvc.Service) *AccountHandler {
	return &AccountHandler{svc: svc, users: users}
}

// ForgotPassword always answers 202 so the response does not reveal whether
// the address has an account.
func (h *AccountHandler) ForgotPassword(c *gin.Context) {
	var req accountdto.ForgotPasswordRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if err := h.svc.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *AccountHandler) ResetPassword(c *gin.Context) {
	var req accountdto.ResetPasswordRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AccountHandler) VerifyEmail(c *gin.Context) {
	var req accountdto.TokenRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if err := h.svc.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ResendVerification emails a new verification link to the signed-in user.
func (h *AccountHandler) ResendVerification(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	if err := h.svc.SendVerification(c.Request.Context(), userID); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *AccountHandler) InspectInvite(c *gin.Context) {
	var req accountdto.TokenRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	invitation, err := h.svc.InspectInvite(c.Request.Context(), req.Token)
	httpx.OK(c, invitation, err)
}

func (h *AccountHandler) AcceptInvite(c *gin.Context) {
	var req accountdto.AcceptInviteRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	if err := h.svc.AcceptInvite(c.Request.Context(), req.Token, req.Password, req.Name); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// InviteUser creates an invited account and emails the invitation (admin).
func (h *AccountHandler) InviteUser(c *gin.Context) {
	var req accountdto.InviteUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}
	user, err := h.svc.Invite(c.Request.Context(), accountsvc.InviteRequest{
		Email: req.Email, FullName: req.Name, RoleID: req.RoleID,
	}, middleware.GetUserEmail(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	item, err := h.users.ResponseForUser(c.Request.Context(), user)
	httpx.Created(c, item, err)
}

// ResendInvite emails a fresh invitation link (admin).
func (h *AccountHandler) ResendInvite(c *gin.Context) {
	if err := h.svc.ResendInvite(c.Request.Context(), c.Param("id"), middleware.GetUserEmail(c)); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}
