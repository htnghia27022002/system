package media_test

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"testing"

	"be/internal/services/media"
)

func TestSaveAvatarRejectsTooLarge(t *testing.T) {
	svc := media.NewService(t.TempDir())
	payload := bytes.Repeat([]byte("a"), media.MaxAvatarBytes+10)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="big.png"`)
	h.Set("Content-Type", "image/png")
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(payload)
	_ = w.Close()

	reader := multipart.NewReader(&buf, w.Boundary())
	form, err := reader.ReadForm(int64(len(payload)) + 1024)
	if err != nil {
		t.Fatal(err)
	}
	defer form.RemoveAll()
	f, err := form.File["file"][0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = svc.SaveAvatar(f, form.File["file"][0])
	if err == nil || !strings.Contains(err.Error(), "2 MB") {
		t.Fatalf("expected size error, got %v", err)
	}
}

func TestResolveAvatarFileBlocksTraversal(t *testing.T) {
	svc := media.NewService(t.TempDir())
	if path := svc.ResolveAvatarFile("../etc/passwd"); path != "" {
		t.Fatalf("expected empty for traversal, got %s", path)
	}
	if path := svc.ResolveAvatarFile("not-uuid.txt"); path != "" {
		t.Fatalf("expected empty for bad ext, got %s", path)
	}
}

func TestSaveAvatarWritesPNG(t *testing.T) {
	dir := t.TempDir()
	svc := media.NewService(dir)
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
	part, _ := w.CreatePart(h)
	_, _ = part.Write(png)
	_ = w.Close()

	reader := multipart.NewReader(&buf, w.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatal(err)
	}
	defer form.RemoveAll()
	f, err := form.File["file"][0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	publicPath, err := svc.SaveAvatar(f, form.File["file"][0])
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if !strings.HasPrefix(publicPath, "/api/media/avatars/") {
		t.Fatalf("unexpected path %s", publicPath)
	}
	disk := svc.ResolveAvatarFile(publicPath)
	if disk == "" {
		t.Fatal("resolve failed")
	}
	if _, err := os.Stat(disk); err != nil {
		t.Fatalf("file missing: %v", err)
	}
}
