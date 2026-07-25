package user

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	apperrors "be/internal/common/errors"
	usermodel "be/internal/models/user"
)

const (
	MaxPhoneLen       = 50
	MaxGeneralLen     = 1000
	MaxAddressLen     = 500
	MaxSocialLinks    = 5
	MaxSocialLabelLen = 50
	MinNameLen        = 2
)

// SocialLinkDTO is the API shape for one social link entry.
type SocialLinkDTO struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// NormalizedPersonal is the result of validating optional personal fields.
// Each Set* flag indicates the field was present in the request (including clear).
type NormalizedPersonal struct {
	PhoneSet    bool
	Phone       string
	GeneralSet  bool
	General     string
	AddressSet  bool
	Address     string
	BirthdaySet bool
	Birthday    *time.Time // nil means clear when BirthdaySet
	LinksSet    bool
	Links       []usermodel.SocialLink
}

// NormalizeOptionalString trims; whitespace-only becomes empty.
func NormalizeOptionalString(v string) string {
	return strings.TrimSpace(v)
}

// ValidatePersonalFields validates optional personal fields.
// Nil pointers mean "not provided" (leave unchanged on apply).
func ValidatePersonalFields(phone, general, birthday, address *string, socialLinks *[]SocialLinkDTO) (*NormalizedPersonal, error) {
	out := &NormalizedPersonal{}

	if phone != nil {
		p := NormalizeOptionalString(*phone)
		if len(p) > MaxPhoneLen {
			return nil, fmt.Errorf("%w: phone must be at most %d characters", apperrors.ErrBadRequest, MaxPhoneLen)
		}
		out.PhoneSet = true
		out.Phone = p
	}
	if general != nil {
		g := NormalizeOptionalString(*general)
		if len(g) > MaxGeneralLen {
			return nil, fmt.Errorf("%w: general must be at most %d characters", apperrors.ErrBadRequest, MaxGeneralLen)
		}
		out.GeneralSet = true
		out.General = g
	}
	if address != nil {
		a := NormalizeOptionalString(*address)
		if len(a) > MaxAddressLen {
			return nil, fmt.Errorf("%w: address must be at most %d characters", apperrors.ErrBadRequest, MaxAddressLen)
		}
		out.AddressSet = true
		out.Address = a
	}
	if birthday != nil {
		raw := strings.TrimSpace(*birthday)
		out.BirthdaySet = true
		if raw != "" {
			parsed, parseErr := time.Parse("2006-01-02", raw)
			if parseErr != nil {
				return nil, fmt.Errorf("%w: birthday must be YYYY-MM-DD", apperrors.ErrBadRequest)
			}
			if raw > time.Now().Format("2006-01-02") {
				return nil, fmt.Errorf("%w: birthday cannot be in the future", apperrors.ErrBadRequest)
			}
			day := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
			out.Birthday = &day
		}
	}
	if socialLinks != nil {
		if len(*socialLinks) > MaxSocialLinks {
			return nil, fmt.Errorf("%w: at most %d social links allowed", apperrors.ErrBadRequest, MaxSocialLinks)
		}
		links := make([]usermodel.SocialLink, 0, len(*socialLinks))
		for i, link := range *socialLinks {
			label := NormalizeOptionalString(link.Label)
			if len(label) > MaxSocialLabelLen {
				return nil, fmt.Errorf("%w: social link label must be at most %d characters", apperrors.ErrBadRequest, MaxSocialLabelLen)
			}
			u := NormalizeOptionalString(link.URL)
			if u == "" {
				return nil, fmt.Errorf("%w: social link %d url is required", apperrors.ErrBadRequest, i+1)
			}
			parsed, parseErr := url.ParseRequestURI(u)
			if parseErr != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
				return nil, fmt.Errorf("%w: social link %d url must be http or https", apperrors.ErrBadRequest, i+1)
			}
			links = append(links, usermodel.SocialLink{Label: label, URL: u})
		}
		out.LinksSet = true
		out.Links = links
	}
	return out, nil
}

// Apply writes normalized optional fields onto the user model.
func (n *NormalizedPersonal) Apply(user *usermodel.User) {
	if n == nil {
		return
	}
	if n.PhoneSet {
		user.Phone = n.Phone
	}
	if n.GeneralSet {
		user.General = n.General
	}
	if n.AddressSet {
		user.Address = n.Address
	}
	if n.BirthdaySet {
		user.Birthday = n.Birthday
	}
	if n.LinksSet {
		user.SocialLinks = n.Links
		if user.SocialLinks == nil {
			user.SocialLinks = []usermodel.SocialLink{}
		}
	}
}

// FormatBirthday returns YYYY-MM-DD or nil.
func FormatBirthday(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format("2006-01-02")
	return &s
}

// SocialLinksToDTO converts model links to DTO slice (never nil).
func SocialLinksToDTO(links []usermodel.SocialLink) []SocialLinkDTO {
	if len(links) == 0 {
		return []SocialLinkDTO{}
	}
	out := make([]SocialLinkDTO, 0, len(links))
	for _, l := range links {
		out = append(out, SocialLinkDTO{Label: l.Label, URL: l.URL})
	}
	return out
}

// ValidateName checks required display name.
func ValidateName(name string) error {
	n := strings.TrimSpace(name)
	if len(n) < MinNameLen {
		return fmt.Errorf("%w: name must be at least %d characters", apperrors.ErrBadRequest, MinNameLen)
	}
	return nil
}
