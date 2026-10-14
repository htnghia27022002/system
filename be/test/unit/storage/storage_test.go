package storage_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"be/pkg/storage"
)

func TestCleanKey(t *testing.T) {
	good := []string{"avatars/a.png", "a", "module/x/y.pdf"}
	bad := []string{"", "/abs", "../x", "a/../../x", "a//b", "a/./b", ".hidden", "a/.upload-1", `a\b`, "a/"}
	for _, k := range good {
		if _, err := storage.CleanKey(k); err != nil {
			t.Errorf("%q should be valid: %v", k, err)
		}
	}
	for _, k := range bad {
		if _, err := storage.CleanKey(k); err == nil {
			t.Errorf("%q should be rejected", k)
		}
	}
}

func TestLocalRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := storage.NewLocal(t.TempDir(), "https://api.example.com/api/media/")
	if err := s.Put(ctx, "avatars/x.png", strings.NewReader("png-bytes"), 9, "image/png"); err != nil {
		t.Fatal(err)
	}
	u := s.URL("avatars/x.png")
	if u != "https://api.example.com/api/media/avatars/x.png" {
		t.Fatalf("url %q", u)
	}
	if key, ok := s.KeyFromURL(u); !ok || key != "avatars/x.png" {
		t.Fatalf("KeyFromURL %q %v", key, ok)
	}
	if _, ok := s.KeyFromURL("https://evil.example.com/api/media/avatars/x.png"); ok {
		t.Fatal("foreign URL must not map to a key")
	}
	rc, obj, err := s.Open(ctx, "avatars/x.png")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(rc)
	rc.Close()
	if string(body) != "png-bytes" || obj.Size != 9 || obj.ContentType != "image/png" {
		t.Fatalf("open: %q %+v", body, obj)
	}
	if err := s.Delete(ctx, "avatars/x.png"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Open(ctx, "avatars/x.png"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if _, _, err := s.Open(ctx, "../../etc/passwd"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("traversal must read as not found, got %v", err)
	}
}

// fakeS3 is a path-style, in-memory S3 that requires SigV4-signed requests.
type fakeS3 struct {
	mu      sync.Mutex
	objects map[string][]byte
	types   map[string]string
}

func (f *fakeS3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=AKIDTEST/") {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "bucket" || path == "bucket/" { // BucketExists
		w.WriteHeader(http.StatusOK)
		return
	}
	key := strings.TrimPrefix(path, "bucket/")
	switch r.Method {
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		if strings.HasPrefix(r.Header.Get("X-Amz-Content-Sha256"), "STREAMING-") {
			body = decodeAWSChunked(body)
		}
		f.objects[key] = body
		f.types[key] = r.Header.Get("Content-Type")
		w.Header().Set("ETag", `"etag"`)
	case http.MethodGet, http.MethodHead:
		body, ok := f.objects[key]
		if !ok {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusNotFound)
			if r.Method == http.MethodGet {
				_, _ = io.WriteString(w, `<Error><Code>NoSuchKey</Code><Message>missing</Message></Error>`)
			}
			return
		}
		w.Header().Set("Content-Type", f.types[key])
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Header().Set("ETag", `"etag"`)
		w.Header().Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
		if r.Method == http.MethodGet {
			_, _ = w.Write(body)
		}
	case http.MethodDelete:
		delete(f.objects, key)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// decodeAWSChunked strips SigV4 streaming framing:
// "<hex-size>;chunk-signature=<sig>\r\n<data>\r\n" ... "0;chunk-signature=…\r\n\r\n".
func decodeAWSChunked(b []byte) []byte {
	var out []byte
	for len(b) > 0 {
		line, rest, ok := bytes.Cut(b, []byte("\r\n"))
		if !ok {
			break
		}
		sizeHex, _, _ := bytes.Cut(line, []byte(";"))
		size, err := strconv.ParseInt(string(sizeHex), 16, 64)
		if err != nil || size == 0 || int64(len(rest)) < size {
			break
		}
		out = append(out, rest[:size]...)
		b = bytes.TrimPrefix(rest[size:], []byte("\r\n"))
	}
	return out
}

func TestS3RoundTripAgainstFakeServer(t *testing.T) {
	fake := &fakeS3{objects: map[string][]byte{}, types: map[string]string{}}
	srv := httptest.NewServer(fake)
	defer srv.Close()

	s, err := storage.NewS3(storage.S3Config{
		Endpoint: strings.TrimPrefix(srv.URL, "http://"), Region: "us-east-1", Bucket: "bucket",
		AccessKeyID: "AKIDTEST", SecretAccessKey: "secret", UseSSL: false, ForcePathStyle: true,
		PublicURL: "https://cdn.example.com/",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	data := []byte("hello avatar")
	if err := s.Put(ctx, "avatars/a b.png", bytes.NewReader(data), int64(len(data)), "image/png"); err != nil {
		t.Fatalf("put: %v", err)
	}
	if got := string(fake.objects["avatars/a b.png"]); got != "hello avatar" || fake.types["avatars/a b.png"] != "image/png" {
		t.Fatalf("server got %q (%s)", got, fake.types["avatars/a b.png"])
	}
	u := s.URL("avatars/a b.png")
	if u != "https://cdn.example.com/avatars/a%20b.png" {
		t.Fatalf("url %q", u)
	}
	rc, obj, err := s.Open(ctx, "avatars/a b.png")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	body, _ := io.ReadAll(rc)
	rc.Close()
	if string(body) != "hello avatar" || obj.ContentType != "image/png" {
		t.Fatalf("open: %q %+v", body, obj)
	}
	if err := s.Delete(ctx, "avatars/a b.png"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, _, err := s.Open(ctx, "avatars/a b.png"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
	if _, err := storage.NewS3(storage.S3Config{Endpoint: "x"}); err == nil {
		t.Fatal("missing bucket/public URL must fail fast")
	}
}
