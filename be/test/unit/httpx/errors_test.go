package httpx_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	apperrors "be/common/errors"
	"be/common/httpx"
	"be/common/response"
)

type signup struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func call(t *testing.T, h gin.HandlerFunc, body string) (int, map[string]any) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/x", h)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(body)))
	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("body is not JSON: %s", w.Body)
	}
	return w.Code, out
}

func TestValidationErrorsUseJSONFieldNames(t *testing.T) {
	code, body := call(t, func(c *gin.Context) {
		var req signup
		if httpx.BindJSON(c, &req) {
			c.Status(http.StatusNoContent)
		}
	}, `{"email":"nope","password":"short"}`)
	fields, _ := body["fields"].(map[string]any)
	if code != http.StatusBadRequest || body["code"] != "validation_failed" ||
		fields["email"] != "email" || fields["password"] != "min=8" {
		t.Fatalf("got %d %v", code, body)
	}
}

func TestMalformedJSONHasItsOwnCode(t *testing.T) {
	for _, raw := range []string{`{"email":`, ``, `{"email": 5}`} {
		code, body := call(t, func(c *gin.Context) {
			var req signup
			httpx.BindJSON(c, &req)
		}, raw)
		if code != http.StatusBadRequest || body["code"] != "invalid_json" {
			t.Fatalf("%q: got %d %v", raw, code, body)
		}
	}
}

func TestHandleErrorAddsDefaultOrSpecificCode(t *testing.T) {
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{apperrors.ErrNotFound, 404, "not_found"},
		{fmt.Errorf("%w: nope", apperrors.ErrForbidden), 403, "forbidden"},
		{apperrors.WithCode(fmt.Errorf("%w: verify", apperrors.ErrForbidden), "email_not_verified"), 403, "email_not_verified"},
		{fmt.Errorf("pq: secret detail"), 500, "internal"},
	}
	for _, tc := range cases {
		status, body := call(t, func(c *gin.Context) { response.HandleError(c, tc.err) }, "")
		if status != tc.status || body["code"] != tc.code {
			t.Fatalf("%v: got %d %v", tc.err, status, body)
		}
		if status == 500 && strings.Contains(fmt.Sprint(body["error"]), "secret") {
			t.Fatal("internal errors must not leak")
		}
	}
}
