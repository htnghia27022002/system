package user_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	apperrors "be/internal/common/errors"
	userdto "be/internal/dto/user"
)

func TestValidatePersonalFieldsRejectsFutureBirthday(t *testing.T) {
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	_, err := userdto.ValidatePersonalFields(nil, nil, &tomorrow, nil, nil)
	if err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request for future birthday, got %v", err)
	}
}

func TestValidatePersonalFieldsAllowsTodayBirthday(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	n, err := userdto.ValidatePersonalFields(nil, nil, &today, nil, nil)
	if err != nil {
		t.Fatalf("today birthday should be allowed: %v", err)
	}
	if !n.BirthdaySet || n.Birthday == nil {
		t.Fatal("expected birthday set")
	}
}

func TestValidatePersonalFieldsRejectsBadSocialURL(t *testing.T) {
	links := []userdto.SocialLinkDTO{{Label: "x", URL: "ftp://example.com"}}
	_, err := userdto.ValidatePersonalFields(nil, nil, nil, nil, &links)
	if err == nil {
		t.Fatal("expected error for non-http(s) url")
	}
}

func TestValidatePersonalFieldsRejectsTooManyLinks(t *testing.T) {
	links := make([]userdto.SocialLinkDTO, 6)
	for i := range links {
		links[i] = userdto.SocialLinkDTO{URL: "https://example.com/" + string(rune('a'+i))}
	}
	_, err := userdto.ValidatePersonalFields(nil, nil, nil, nil, &links)
	if err == nil {
		t.Fatal("expected error for more than 5 links")
	}
}

func TestValidatePersonalFieldsRejectsOverlongGeneral(t *testing.T) {
	g := strings.Repeat("a", userdto.MaxGeneralLen+1)
	_, err := userdto.ValidatePersonalFields(nil, &g, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for overlong general")
	}
}

func TestValidatePersonalFieldsTrimsWhitespacePhone(t *testing.T) {
	phone := "   "
	n, err := userdto.ValidatePersonalFields(&phone, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if !n.PhoneSet || n.Phone != "" {
		t.Fatalf("expected trimmed empty phone, got %q", n.Phone)
	}
}
