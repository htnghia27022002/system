package user

import "be/pkg/query"

type CreateUserRequest struct {
	Email       string          `json:"email" binding:"required,email"`
	Password    string          `json:"password" binding:"required,min=8"`
	Name        string          `json:"name" binding:"required,min=2"`
	RoleID      string          `json:"roleId" binding:"required,uuid"`
	Status      string          `json:"status" binding:"omitempty,oneof=active inactive"`
	Phone       *string         `json:"phone"`
	General     *string         `json:"general"`
	Birthday    *string         `json:"birthday"`
	Address     *string         `json:"address"`
	SocialLinks *[]SocialLinkDTO `json:"socialLinks"`
	AvatarURL   *string         `json:"avatarUrl"`
}

type UpdateUserRequest struct {
	Email       *string         `json:"email" binding:"omitempty,email"`
	Name        *string         `json:"name" binding:"omitempty,min=2"`
	Password    *string         `json:"password" binding:"omitempty,min=8"`
	RoleID      *string         `json:"roleId" binding:"omitempty,uuid"`
	Status      *string         `json:"status" binding:"omitempty,oneof=active inactive"`
	Phone       *string         `json:"phone"`
	General     *string         `json:"general"`
	Birthday    *string         `json:"birthday"`
	Address     *string         `json:"address"`
	SocialLinks *[]SocialLinkDTO `json:"socialLinks"`
	AvatarURL   *string         `json:"avatarUrl"`
}

type UserResponse struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	RoleID         string          `json:"roleId"`
	Status         string          `json:"status"`
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
	Status string `form:"status" binding:"omitempty,oneof=active inactive"`
}

type PaginatedUsersResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}
