package routes

import (
	"net/http"

	"be/internal/apidoc"
	accountdto "be/internal/dto/account"
	admindto "be/internal/dto/admin"
	authdto "be/internal/dto/auth"
	permissiondto "be/internal/dto/permission"
	roledto "be/internal/dto/role"
	userdto "be/internal/dto/user"
	accountsvc "be/internal/services/account"
	"be/pkg/query"
	"be/public/handlers"
)

// CoreOperations documents every core route (modules document theirs through
// module.APIDocProvider). Keep this next to the route files: the OpenAPI test
// fails when a mounted route is missing here, or listed here but not mounted.
func CoreOperations() []apidoc.Op {
	type op = apidoc.Op
	const (
		app     = "App"
		auth    = "Auth"
		account = "Account"
		users   = "Users"
		roles   = "Roles"
		system  = "System"
	)
	limited := []string{"too_many_requests"}
	return []op{
		// App
		{Method: "GET", Path: "/api/app/info", Tag: app, Summary: "Enabled modules and public settings", Response: handlers.AppInfoResponse{}},
		{Method: "GET", Path: "/api/health", Tag: app, Summary: "Liveness probe", Response: handlers.HealthResponse{}},
		{Method: "GET", Path: "/api/ready", Tag: app, Summary: "Readiness probe (503 when a dependency is down)", Response: handlers.ReadyResponse{}},
		{Method: "GET", Path: "/api/openapi.json", Tag: app, Summary: "This OpenAPI document", Raw: "application/json"},
		{Method: "GET", Path: "/api/media/*key", Tag: app, Summary: "Stream an uploaded file (local storage)", Raw: "application/octet-stream"},

		// Auth
		{Method: "POST", Path: "/api/auth/login", Tag: auth, Summary: "Sign in with email and password; the refresh token is set as an httpOnly cookie",
			Body: authdto.LoginRequest{}, Response: authdto.AuthResponse{}, Codes: []string{"unauthorized", "account_locked", "email_not_verified", "too_many_requests"}},
		{Method: "POST", Path: "/api/auth/register", Tag: auth, Summary: "Create an account",
			Body: authdto.RegisterRequest{}, Response: authdto.AuthResponse{}, Status: http.StatusCreated, Codes: []string{"conflict", "registration_disabled", "email_not_verified", "too_many_requests"}},
		{Method: "POST", Path: "/api/auth/refresh", Tag: auth, Summary: "Rotate the token pair (cookie + X-Requested-With, or refreshToken in the body)",
			Body: authdto.RefreshRequest{}, Response: authdto.TokenPairResponse{}},
		{Method: "POST", Path: "/api/auth/logout", Tag: auth, Summary: "Revoke the refresh token and clear the cookie", Body: authdto.LogoutRequest{}},
		{Method: "GET", Path: "/api/auth/oauth/providers", Tag: auth, Summary: "Configured OAuth providers", Response: authdto.OAuthProvidersResponse{}},
		{Method: "GET", Path: "/api/auth/oauth/:provider/start", Tag: auth, Summary: "Redirect to the provider (query: redirect_uri)", Status: http.StatusFound, Raw: "text/html"},
		{Method: "POST", Path: "/api/auth/oauth/:provider/callback", Tag: auth, Summary: "Exchange the OAuth code for a session",
			Body: authdto.OAuthCallbackRequest{}, Response: authdto.AuthResponse{}, Codes: []string{"registration_disabled", "too_many_requests"}},
		{Method: "GET", Path: "/api/auth/me", Tag: auth, Summary: "The signed-in user", Auth: true, Response: authdto.AuthUserResponse{}},
		{Method: "PATCH", Path: "/api/auth/profile", Tag: auth, Summary: "Update own profile", Auth: true, Body: authdto.UpdateProfileRequest{}, Response: authdto.AuthUserResponse{}},
		{Method: "POST", Path: "/api/auth/profile/avatar", Tag: auth, Summary: "Upload own avatar (multipart field: file)", Auth: true, Response: authdto.AuthUserResponse{}},
		{Method: "POST", Path: "/api/auth/change-password", Tag: auth, Summary: "Change own password", Auth: true, Body: authdto.ChangePasswordRequest{}, Response: authdto.MessageResponse{}, Codes: limited},
		{Method: "GET", Path: "/api/auth/sessions", Tag: auth, Summary: "Own active sessions (X-Session-Id marks the current one)", Auth: true, Response: authdto.SessionListResponse{}},
		{Method: "POST", Path: "/api/auth/sessions/revoke-others", Tag: auth, Summary: "Sign out every other session", Auth: true, Body: authdto.RevokeOthersRequest{}, Response: authdto.SessionListResponse{}},
		{Method: "DELETE", Path: "/api/auth/sessions/:id", Tag: auth, Summary: "Revoke one of own sessions", Auth: true},

		// Account (email links)
		{Method: "POST", Path: "/api/auth/password/forgot", Tag: account, Summary: "Email a reset link (always 202)", Body: accountdto.ForgotPasswordRequest{}, Status: http.StatusAccepted, Codes: limited},
		{Method: "POST", Path: "/api/auth/password/reset", Tag: account, Summary: "Set a new password from a reset link", Body: accountdto.ResetPasswordRequest{}, Codes: []string{"bad_request", "too_many_requests"}},
		{Method: "POST", Path: "/api/auth/email/verify", Tag: account, Summary: "Confirm the email from a verification link", Body: accountdto.TokenRequest{}},
		{Method: "POST", Path: "/api/auth/email/verification", Tag: account, Summary: "Resend the verification email", Auth: true, Status: http.StatusAccepted, Codes: limited},
		{Method: "POST", Path: "/api/auth/invite/inspect", Tag: account, Summary: "Show who an invitation link is for", Body: accountdto.TokenRequest{}, Response: accountsvc.Invitation{}},
		{Method: "POST", Path: "/api/auth/invite/accept", Tag: account, Summary: "Set a password and activate an invited account", Body: accountdto.AcceptInviteRequest{}},

		// Admin: system
		{Method: "GET", Path: "/api/admin/dashboard/overview", Tag: system, Summary: "Dashboard KPIs", Auth: true, Permission: "dashboard:view", Response: admindto.DashboardOverview{}},
		{Method: "GET", Path: "/api/admin/audit-logs", Tag: system, Summary: "Audit trail", Auth: true, Permission: "audit:view", Query: admindto.ListAuditLogsQuery{}, Response: pageOf[admindto.AuditLogResponse]()},
		{Method: "GET", Path: "/api/admin/settings", Tag: system, Summary: "All settings with effective values", Auth: true, Permission: "settings:view", Response: admindto.SettingsResponse{}},
		{Method: "PATCH", Path: "/api/admin/settings", Tag: system, Summary: "Update settings (partial)", Auth: true, Permission: "settings:modify", Body: admindto.UpdateSettingsRequest{}, Response: admindto.SettingsResponse{}},

		// Admin: users
		{Method: "GET", Path: "/api/admin/users", Tag: users, Summary: "List users", Auth: true, Permission: "users:view", Query: userdto.ListUsersQuery{}, Response: userdto.PaginatedUsersResponse{}},
		{Method: "POST", Path: "/api/admin/users", Tag: users, Summary: "Create a user", Auth: true, Permission: "users:modify", Body: userdto.CreateUserRequest{}, Response: userdto.UserResponse{}, Status: http.StatusCreated},
		{Method: "POST", Path: "/api/admin/users/invite", Tag: users, Summary: "Invite a user by email", Auth: true, Permission: "users:modify", Body: accountdto.InviteUserRequest{}, Response: userdto.UserResponse{}, Status: http.StatusCreated},
		{Method: "POST", Path: "/api/admin/users/:id/invite", Tag: users, Summary: "Resend an invitation", Auth: true, Permission: "users:modify", Status: http.StatusAccepted},
		{Method: "GET", Path: "/api/admin/users/:id", Tag: users, Summary: "Get a user", Auth: true, Permission: "users:view", Response: userdto.UserResponse{}},
		{Method: "PATCH", Path: "/api/admin/users/:id", Tag: users, Summary: "Update a user", Auth: true, Permission: "users:modify", Body: userdto.UpdateUserRequest{}, Response: userdto.UserResponse{}},
		{Method: "POST", Path: "/api/admin/users/:id/avatar", Tag: users, Summary: "Upload a user's avatar (multipart field: file)", Auth: true, Permission: "users:modify", Response: userdto.UserResponse{}},
		{Method: "DELETE", Path: "/api/admin/users/:id", Tag: users, Summary: "Delete a user", Auth: true, Permission: "users:modify"},

		// Admin: roles & permissions
		{Method: "GET", Path: "/api/admin/roles", Tag: roles, Summary: "List roles", Auth: true, Permission: "roles:view", Query: roledto.ListRolesQuery{}, Response: roledto.PaginatedRolesResponse{}},
		{Method: "GET", Path: "/api/admin/roles/all", Tag: roles, Summary: "All roles (unpaged)", Auth: true, Permission: "roles:view", Response: []roledto.RoleResponse{}},
		{Method: "POST", Path: "/api/admin/roles", Tag: roles, Summary: "Create a role", Auth: true, Permission: "roles:modify", Body: roledto.CreateRoleRequest{}, Response: roledto.RoleResponse{}, Status: http.StatusCreated},
		{Method: "GET", Path: "/api/admin/roles/:id", Tag: roles, Summary: "Get a role", Auth: true, Permission: "roles:view", Response: roledto.RoleResponse{}},
		{Method: "PATCH", Path: "/api/admin/roles/:id", Tag: roles, Summary: "Update a role", Auth: true, Permission: "roles:modify", Body: roledto.UpdateRoleRequest{}, Response: roledto.RoleResponse{}},
		{Method: "DELETE", Path: "/api/admin/roles/:id", Tag: roles, Summary: "Delete a role", Auth: true, Permission: "roles:modify"},
		{Method: "GET", Path: "/api/admin/permissions", Tag: roles, Summary: "List permissions", Auth: true, Permission: "permissions:view", Query: permissiondto.ListPermissionsQuery{}, Response: permissiondto.PaginatedPermissionsResponse{}},
		{Method: "GET", Path: "/api/admin/permissions/all", Tag: roles, Summary: "All permissions (unpaged)", Auth: true, Permission: "permissions:view", Response: []permissiondto.PermissionResponse{}},
	}
}

func pageOf[T any]() query.Page[T] { return query.Page[T]{} }
