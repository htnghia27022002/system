package oauth

import (
	"fmt"
	"net/url"

	apperrors "be/internal/common/errors"
)

// ValidateRedirectURI checks that a redirect URI is an absolute URL with scheme and host.
func ValidateRedirectURI(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: invalid redirect uri", apperrors.ErrBadRequest)
	}
	return nil
}
