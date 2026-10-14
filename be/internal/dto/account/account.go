package account

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required,max=200"`
	Password string `json:"password" binding:"required,min=8,max=200"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"required,max=200"`
}

type AcceptInviteRequest struct {
	Token    string `json:"token" binding:"required,max=200"`
	Password string `json:"password" binding:"required,min=8,max=200"`
	Name     string `json:"name" binding:"omitempty,min=2,max=255"`
}

type InviteUserRequest struct {
	Email  string `json:"email" binding:"required,email,max=255"`
	Name   string `json:"name" binding:"required,min=2,max=255"`
	RoleID string `json:"roleId" binding:"required,uuid"`
}
