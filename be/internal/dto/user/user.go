package user

import "be/pkg/query"

type CreateUserRequest struct {
	Email       string           `json:"email" binding:"required,email"`
	Password    string           `json:"password" binding:"required,min=8"`
	Name        string           `json:"name" binding:"required,min=2"`
	RoleID      string           `json:"roleId" binding:"required,uuid"`
	Status      string           `json:"status" binding:"omitempty,oneof=active inactive"`
	Phone       *string          `json:"phone"`
	General     *string          `json:"general"`
	Birthday    *string          `json:"birthday"`
	Address     *string          `json:"address"`
	SocialLinks *[]SocialLinkDTO `json:"socialLinks"`
	AvatarURL   *string          `json:"avatarUrl"`
	SuperAdmin  *bool            `json:"superAdmin"`
}

type UpdateUserRequest struct {
	Email       *string          `json:"email" binding:"omitempty,email"`
	Name        *string          `json:"name" binding:"omitempty,min=2"`
	Password    *string          `json:"password" binding:"omitempty,min=8"`
	RoleID      *string          `json:"roleId" binding:"omitempty,uuid"`
	Status      *string          `json:"status" binding:"omitempty,oneof=active inactive"`
	Phone       *string          `json:"phone"`
	General     *string          `json:"general"`
	Birthday    *string          `json:"birthday"`
	Address     *string          `json:"address"`
	SocialLinks *[]SocialLinkDTO `json:"socialLinks"`
	AvatarURL   *string          `json:"avatarUrl"`
	SuperAdmin  *bool            `json:"superAdmin"`
}

type UserResponse struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	RoleID         string          `json:"roleId"`
	SuperAdmin     bool            `json:"superAdmin"`
	Status         string          `json:"status"`
	EmailVerified  bool            `json:"emailVerified"`
	Phone          string          `json:"phone"`
	AvatarURL      string          `json:"avatarUrl"`
	General        string          `json:"general"`
	Birthday       *string         `json:"birthday"`
	Address        string          `json:"address"`
	SocialLinks    []SocialLinkDTO `json:"socialLinks"`
	OAuthProviders []string        `json:"oauthProviders"`
	CreatedAt      string          `json:"createdAt"`
}

type ListUsersQuery struct {
	query.PageParams
	Search string `form:"search" binding:"omitempty,max=200"`
	RoleID string `form:"roleId" binding:"omitempty,uuid"`
	ID     string `form:"id" binding:"omitempty,uuid"`
	Status string `form:"status" binding:"omitempty,oneof=active inactive invited"`
}

type PaginatedUsersResponse = query.Page[UserResponse]

// ChangedFields lists the JSON names of the fields an update sets (never their
// values — the audit log must not hold passwords or personal data).
func ChangedFields(req UpdateUserRequest) []string {
	fields := []string{}
	add := func(set bool, name string) {
		if set {
			fields = append(fields, name)
		}
	}
	add(req.Email != nil, "email")
	add(req.Name != nil, "name")
	add(req.Password != nil && *req.Password != "", "password")
	add(req.RoleID != nil, "roleId")
	add(req.Status != nil, "status")
	add(req.Phone != nil, "phone")
	add(req.General != nil, "general")
	add(req.Birthday != nil, "birthday")
	add(req.Address != nil, "address")
	add(req.SocialLinks != nil, "socialLinks")
	add(req.AvatarURL != nil, "avatarUrl")
	add(req.SuperAdmin != nil, "superAdmin")
	return fields
}
