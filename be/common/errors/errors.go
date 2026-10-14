package errors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
	// ErrTooManyRequests maps to 429 (rate limit, temporary account lockout).
	ErrTooManyRequests = errors.New("too many requests")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// CodedError attaches a stable, machine-readable code (e.g. "email_not_verified")
// to an error so clients can branch or translate without parsing the message.
// errors.Is still sees the wrapped sentinel, so the HTTP status is unchanged.
type CodedError struct {
	Code string
	Err  error
}

func (e *CodedError) Error() string { return e.Err.Error() }
func (e *CodedError) Unwrap() error { return e.Err }

// WithCode wraps err with a client-facing code.
func WithCode(err error, code string) error {
	return &CodedError{Code: code, Err: err}
}

// Code returns the first code in err's chain, or "".
func Code(err error) string {
	var coded *CodedError
	if errors.As(err, &coded) {
		return coded.Code
	}
	return ""
}
