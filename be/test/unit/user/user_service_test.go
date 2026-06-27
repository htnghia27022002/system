package user_test

import (
	"context"
	"testing"

	apperrors "be/internal/common/errors"
	userdto "be/internal/dto/user"
	usermodel "be/internal/models/user"
	usersvc "be/internal/services/user"
	"be/test/testutil"
)

func TestCreateHashesPassword(t *testing.T) {
	repo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{}}
	svc := usersvc.NewService(repo)

	user, err := svc.Create(context.Background(), userdto.CreateUserRequest{
		Email:    "demo@example.com",
		Password: "password123",
		Name:     "Demo User",
		RoleID:   "22222222-2222-4222-8222-222222222222",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.PasswordHash == "password123" {
		t.Fatal("expected hashed password")
	}
}

func TestDeleteBlocksSelfDelete(t *testing.T) {
	repo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"self-id": {ID: "self-id", Email: "self@example.com", FullName: "Self"},
		},
	}
	svc := usersvc.NewService(repo)

	err := svc.Delete(context.Background(), "self-id", "self-id")
	if err == nil || !apperrors.IsForbidden(err) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}
