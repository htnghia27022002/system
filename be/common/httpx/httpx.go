package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	apperrors "be/common/errors"
	"be/common/response"
	"be/internal/middleware"
)

// BindJSON binds the request body. Writes 400 and returns false on failure.
func BindJSON(c *gin.Context, dest any) bool {
	if err := c.ShouldBindJSON(dest); err != nil {
		writeBindError(c, err)
		return false
	}
	return true
}

// BindQuery binds URL query parameters. Writes 400 and returns false on failure.
func BindQuery(c *gin.Context, dest any) bool {
	if err := c.ShouldBindQuery(dest); err != nil {
		writeBindError(c, err)
		return false
	}
	return true
}

func init() {
	// Report validation errors by JSON / form field name, not Go field name.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			for _, tag := range []string{"json", "form"} {
				name := strings.SplitN(f.Tag.Get(tag), ",", 2)[0]
				if name == "-" {
					return ""
				}
				if name != "" {
					return name
				}
			}
			return f.Name
		})
	}
}

// writeBindError turns binding failures into a stable 400 body: field errors
// as validation_failed with {field: rule}, malformed bodies as invalid_json.
func writeBindError(c *gin.Context, err error) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fields := make(map[string]string, len(verrs))
		parts := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			name := fe.Field()
			rule := fe.Tag()
			if fe.Param() != "" {
				rule += "=" + fe.Param()
			}
			fields[name] = rule
			parts = append(parts, name+": "+rule)
		}
		response.ValidationError(c, "invalid request: "+strings.Join(parts, ", "), fields)
		return
	}
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		response.JSON(c, http.StatusBadRequest, gin.H{"error": "request body must be valid JSON", "code": response.CodeInvalidJSON})
		return
	}
	response.Error(c, http.StatusBadRequest, err.Error())
}

// RequireUserID returns the JWT user id or writes 401 and false.
func RequireUserID(c *gin.Context) (string, bool) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.HandleError(c, apperrors.ErrUnauthorized)
		return "", false
	}
	return userID, true
}

// OK writes payload at 200, or maps err through HandleError.
func OK(c *gin.Context, payload any, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, payload)
}

// Created writes payload at 201, or maps err through HandleError.
func Created(c *gin.Context, payload any, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, payload)
}

// NoContent writes 204, or maps err through HandleError.
func NoContent(c *gin.Context, err error) {
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// FormFile opens the first matching multipart field name.
func FormFile(c *gin.Context, errMessage string, names ...string) (multipart.File, *multipart.FileHeader, bool) {
	for _, name := range names {
		file, header, err := c.Request.FormFile(name)
		if err == nil {
			return file, header, true
		}
	}
	if strings.TrimSpace(errMessage) == "" {
		errMessage = "file is required"
	}
	response.Error(c, http.StatusBadRequest, errMessage)
	return nil, nil, false
}

// ValuesMap converts url.Values / headers into a JSON-friendly map.
func ValuesMap(values map[string][]string) map[string]any {
	out := make(map[string]any, len(values))
	for key, items := range values {
		if len(items) == 1 {
			out[key] = items[0]
			continue
		}
		copied := make([]string, len(items))
		copy(copied, items)
		out[key] = copied
	}
	return out
}

// HeaderMap copies an HTTP header into a JSON-friendly map.
func HeaderMap(h http.Header) map[string]any {
	return ValuesMap(h)
}

// QueryMap copies URL query values into a JSON-friendly map.
func QueryMap(v url.Values) map[string]any {
	return ValuesMap(v)
}

// ClientIP returns the client address. X-Forwarded-For / X-Real-IP are honored only
// when the direct peer is a trusted proxy (engine.SetTrustedProxies, TRUSTED_PROXIES);
// otherwise the socket address is used, so clients cannot spoof their IP.
func ClientIP(c *gin.Context) string {
	return c.ClientIP()
}
