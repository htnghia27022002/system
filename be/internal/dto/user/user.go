package user

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
	RoleID   string `json:"roleId" binding:"required,uuid"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Name     *string `json:"name" binding:"omitempty,min=2"`
	Password *string `json:"password" binding:"omitempty,min=8"`
	RoleID   *string `json:"roleId" binding:"omitempty,uuid"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	RoleID string `json:"roleId"`
	Status string `json:"status"`
}

type ListUsersQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Search   string `form:"search"`
	RoleID   string `form:"roleId"`
}

type PaginatedUsersResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}
