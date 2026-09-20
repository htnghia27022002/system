package user_test

import (
	"context"
	"testing"
	"time"

	apperrors "be/common/errors"
	"be/pkg/hash"
	userdto "be/internal/dto/user"
	usermodel "be/internal/models/user"
	usersvc "be/internal/services/user"
	"be/test/testutil"
)

func TestCreateHashesPassword(t *testing.T) {
	repo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{}}
	svc := usersvc.NewService(repo, nil, nil, nil)

	user, err := svc.Create(context.Background(), userdto.CreateUserRequest{
		Email:    "demo@example.com",
		Password: "password123",
		Name:     "Demo User",
		RoleID:   "22222222-2222-4222-8222-222222222222",
	}, usersvc.Actor{ID: "actor"})
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
	svc := usersvc.NewService(repo, nil, nil, nil)

	err := svc.Delete(context.Background(), "self-id", usersvc.Actor{ID: "self-id"})
	if err == nil || !apperrors.IsForbidden(err) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestUpdateKeepsPasswordWhenBlank(t *testing.T) {
	hashed, err := hash.HashPassword("keep-me-12")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	repo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"u1": {
				ID:           "u1",
				Email:        "a@example.com",
				FullName:     "A",
				PasswordHash: hashed,
				Status:       usermodel.StatusActive,
				SocialLinks:  []usermodel.SocialLink{},
			},
		},
	}
	svc := usersvc.NewService(repo, nil, nil, nil)
	blank := ""
	phone := "555"
	user, err := svc.Update(context.Background(), "u1", userdto.UpdateUserRequest{
		Password: &blank,
		Phone:    &phone,
	}, usersvc.Actor{ID: "admin"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if user.PasswordHash != hashed {
		t.Fatal("blank password must keep existing hash")
	}
	if user.Phone != "555" {
		t.Fatalf("expected phone update, got %q", user.Phone)
	}
}

func TestToResponseIncludesPersonalFields(t *testing.T) {
	bday := time.Date(1991, 2, 3, 0, 0, 0, 0, time.UTC)
	user := &usermodel.User{
		ID:          "u1",
		Email:       "a@example.com",
		FullName:    "Name",
		RoleID:      "r1",
		Status:      usermodel.StatusActive,
		Phone:       "1",
		AvatarURL:   "/api/media/avatars/x.png",
		General:     "g",
		Birthday:    &bday,
		Address:     "addr",
		SocialLinks: []usermodel.SocialLink{{Label: "L", URL: "https://example.com"}},
	}
	resp := usersvc.ToResponse(user)
	if resp.Phone != "1" || resp.AvatarURL == "" || resp.General != "g" || resp.Address != "addr" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Birthday == nil || *resp.Birthday != "1991-02-03" {
		t.Fatalf("unexpected birthday: %v", resp.Birthday)
	}
	if len(resp.SocialLinks) != 1 {
		t.Fatalf("expected social links, got %+v", resp.SocialLinks)
	}
	if resp.SuperAdmin {
		t.Fatal("expected superAdmin false")
	}
}

func TestCreateRejectsSuperAdminFlagFromNonSuperAdmin(t *testing.T) {
	repo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{}}
	svc := usersvc.NewService(repo, nil, nil, nil)
	flag := true
	_, err := svc.Create(context.Background(), userdto.CreateUserRequest{
		Email:      "x@example.com",
		Password:   "password123",
		Name:       "X",
		RoleID:     "22222222-2222-4222-8222-222222222222",
		SuperAdmin: &flag,
	}, usersvc.Actor{ID: "actor"})
	if err == nil || !apperrors.IsForbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestUpdateBlocksLastSuperAdminDemotion(t *testing.T) {
	repo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"sa": {ID: "sa", Email: "sa@example.com", FullName: "SA", IsSuperAdmin: true, Status: usermodel.StatusActive},
		},
	}
	svc := usersvc.NewService(repo, nil, nil, nil)
	flag := false
	_, err := svc.Update(context.Background(), "sa", userdto.UpdateUserRequest{
		SuperAdmin: &flag,
	}, usersvc.Actor{ID: "sa", SuperAdmin: true})
	if err == nil || !apperrors.IsForbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestDeleteBlocksLastSuperAdmin(t *testing.T) {
	repo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"sa": {ID: "sa", Email: "sa@example.com", FullName: "SA", IsSuperAdmin: true},
			"u2": {ID: "u2", Email: "u2@example.com", FullName: "U2"},
		},
	}
	svc := usersvc.NewService(repo, nil, nil, nil)
	err := svc.Delete(context.Background(), "sa", usersvc.Actor{ID: "u2", SuperAdmin: true})
	if err == nil || !apperrors.IsForbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestSuperAdminCanGrantFlag(t *testing.T) {
	repo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"sa": {ID: "sa", Email: "sa@example.com", FullName: "SA", IsSuperAdmin: true, Status: usermodel.StatusActive},
			"u2": {ID: "u2", Email: "u2@example.com", FullName: "U2", Status: usermodel.StatusActive},
		},
	}
	svc := usersvc.NewService(repo, nil, nil, nil)
	flag := true
	user, err := svc.Update(context.Background(), "u2", userdto.UpdateUserRequest{
		SuperAdmin: &flag,
	}, usersvc.Actor{ID: "sa", SuperAdmin: true})
	if err != nil {
		t.Fatalf("grant: %v", err)
	}
	if !user.IsSuperAdmin {
		t.Fatal("expected super admin granted")
	}
}
