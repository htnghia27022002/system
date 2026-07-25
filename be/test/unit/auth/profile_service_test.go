package auth_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/textproto"
	"testing"
	"time"

	"be/pkg/hash"
	apperrors "be/internal/common/errors"
	jwtmanager "be/internal/common/jwt"
	authdto "be/internal/dto/auth"
	userdto "be/internal/dto/user"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	authsvc "be/internal/services/auth"
	"be/internal/services/media"
	"be/test/testutil"
)

func TestMeIncludesProfileFieldsAndHasPassword(t *testing.T) {
	hashed, err := hash.HashPassword("secret123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	roleID := "role-1"
	user := &usermodel.User{
		ID:           "user-1",
		Email:        "u@example.com",
		PasswordHash: hashed,
		FullName:     "User One",
		Phone:        "123",
		General:      "About",
		RoleID:       roleID,
		Status:       usermodel.StatusActive,
		SocialLinks:  []usermodel.SocialLink{},
	}
	authRepo := &testutil.MockAuthRepo{
		Users:         map[string]*usermodel.User{"user-1": user},
		RefreshTokens: map[string]*authmodel.RefreshToken{},
	}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{roleID: {ID: roleID, Slug: "admin", Name: "Admin"}},
		Permissions: map[string][]string{roleID: {}},
	}
	svc := authsvc.NewService(authRepo, &testutil.MockUserRepo{}, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)

	me, err := svc.Me(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if !me.HasPassword {
		t.Fatal("expected hasPassword true")
	}
	if me.Phone != "123" || me.General != "About" {
		t.Fatalf("expected profile fields, got phone=%q general=%q", me.Phone, me.General)
	}
	if me.SocialLinks == nil {
		t.Fatal("socialLinks must not be nil")
	}
}

func TestMeHasPasswordFalseWhenEmptyHash(t *testing.T) {
	roleID := "role-1"
	user := &usermodel.User{
		ID:       "oauth-1",
		Email:    "o@example.com",
		FullName: "OAuth User",
		RoleID:   roleID,
		Status:   usermodel.StatusActive,
	}
	authRepo := &testutil.MockAuthRepo{
		Users:         map[string]*usermodel.User{"oauth-1": user},
		RefreshTokens: map[string]*authmodel.RefreshToken{},
	}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{roleID: {ID: roleID, Slug: "user", Name: "User"}},
		Permissions: map[string][]string{roleID: {}},
	}
	svc := authsvc.NewService(authRepo, &testutil.MockUserRepo{}, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)

	me, err := svc.Me(context.Background(), "oauth-1")
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if me.HasPassword {
		t.Fatal("expected hasPassword false")
	}
}

func TestUpdateProfilePersistsPersonalFields(t *testing.T) {
	roleID := "role-1"
	user := &usermodel.User{
		ID:          "user-1",
		Email:       "u@example.com",
		FullName:    "Old Name",
		RoleID:      roleID,
		Status:      usermodel.StatusActive,
		SocialLinks: []usermodel.SocialLink{},
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{roleID: {ID: roleID, Slug: "admin", Name: "Admin"}},
		Permissions: map[string][]string{roleID: {}},
	}
	svc := authsvc.NewService(&testutil.MockAuthRepo{Users: map[string]*usermodel.User{}}, userRepo, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)

	phone := "+84123"
	general := "Bio"
	birthday := "1990-01-15"
	address := "123 Street"
	links := []userdto.SocialLinkDTO{{Label: "GitHub", URL: "https://github.com/example"}}
	out, err := svc.UpdateProfile(context.Background(), "user-1", authdto.UpdateProfileRequest{
		Name:        "New Name",
		Phone:       &phone,
		General:     &general,
		Birthday:    &birthday,
		Address:     &address,
		SocialLinks: &links,
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if out.Name != "New Name" || out.Phone != phone || out.General != general {
		t.Fatalf("unexpected response: %+v", out)
	}
	if out.Birthday == nil || *out.Birthday != birthday {
		t.Fatalf("expected birthday %s, got %v", birthday, out.Birthday)
	}
	if len(out.SocialLinks) != 1 || out.SocialLinks[0].URL != "https://github.com/example" {
		t.Fatalf("unexpected social links: %+v", out.SocialLinks)
	}
}

func TestChangePasswordSuccessAndRejectWrongCurrent(t *testing.T) {
	hashed, err := hash.HashPassword("oldpass12")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	user := &usermodel.User{
		ID:           "user-1",
		Email:        "u@example.com",
		PasswordHash: hashed,
		FullName:     "User",
		Status:       usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	svc := authsvc.NewService(&testutil.MockAuthRepo{}, userRepo, &testutil.MockRoleRepo{}, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)

	err = svc.ChangePassword(context.Background(), "user-1", "wrong", "newpass12")
	if err == nil || !errors.Is(err, apperrors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized for wrong current, got %v", err)
	}

	if err := svc.ChangePassword(context.Background(), "user-1", "oldpass12", "newpass12"); err != nil {
		t.Fatalf("change password: %v", err)
	}
	if err := hash.ComparePassword(userRepo.Users["user-1"].PasswordHash, "newpass12"); err != nil {
		t.Fatal("expected new password to match")
	}
}

func TestChangePasswordRejectedForOAuthOnly(t *testing.T) {
	user := &usermodel.User{
		ID:       "oauth-1",
		Email:    "o@example.com",
		FullName: "OAuth",
		Status:   usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"oauth-1": user}}
	svc := authsvc.NewService(&testutil.MockAuthRepo{}, userRepo, &testutil.MockRoleRepo{}, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)

	err := svc.ChangePassword(context.Background(), "oauth-1", "anything", "newpass12")
	if err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request, got %v", err)
	}
}

func TestUploadAvatarRejectsNonImage(t *testing.T) {
	dir := t.TempDir()
	mediaSvc := media.NewService(dir)
	user := &usermodel.User{
		ID:       "user-1",
		Email:    "u@example.com",
		FullName: "User",
		Status:   usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	svc := authsvc.NewService(&testutil.MockAuthRepo{}, userRepo, &testutil.MockRoleRepo{}, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, mediaSvc)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="x.txt"`)
	h.Set("Content-Type", "text/plain")
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	_, _ = part.Write([]byte("not an image"))
	_ = w.Close()

	reader := multipart.NewReader(&buf, w.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	defer form.RemoveAll()
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("expected file")
	}
	f, err := files[0].Open()
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	_, err = svc.UploadAvatar(context.Background(), "user-1", f, files[0])
	if err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request for non-image, got %v", err)
	}
	if userRepo.Users["user-1"].AvatarURL != "" {
		t.Fatal("avatar_url must remain empty on failure")
	}
}

func TestUploadAvatarAcceptsPNG(t *testing.T) {
	dir := t.TempDir()
	mediaSvc := media.NewService(dir)
	user := &usermodel.User{
		ID:       "user-1",
		Email:    "u@example.com",
		FullName: "User",
		Status:   usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{},
		Permissions: map[string][]string{},
	}
	svc := authsvc.NewService(&testutil.MockAuthRepo{}, userRepo, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, mediaSvc)

	// Minimal 1x1 PNG
	png := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00,
		0x0c, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x01, 0x00, 0x05, 0xfe, 0xd4, 0xef, 0x00, 0x00,
		0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="a.png"`)
	h.Set("Content-Type", "image/png")
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	_, _ = part.Write(png)
	_ = w.Close()

	reader := multipart.NewReader(&buf, w.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	defer form.RemoveAll()
	files := form.File["file"]
	f, err := files[0].Open()
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	out, err := svc.UploadAvatar(context.Background(), "user-1", f, files[0])
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if out.AvatarURL == "" || out.AvatarURL[:len("/api/media/avatars/")] != "/api/media/avatars/" {
		t.Fatalf("unexpected avatarUrl: %s", out.AvatarURL)
	}
}
