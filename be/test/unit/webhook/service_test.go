package webhook_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	apperrors "be/internal/common/errors"
	webhookdto "be/internal/dto/webhook"
	webhookmodel "be/internal/models/webhook"
	webhooksvc "be/internal/services/webhook"
	"be/test/testutil"
)

func newWebhookService() (*webhooksvc.Service, *testutil.MemoryWebhookInboxRepo, *testutil.MemoryWebhookRequestRepo) {
	inboxes := testutil.NewMemoryWebhookInboxRepo()
	requests := testutil.NewMemoryWebhookRequestRepo()
	return webhooksvc.NewService(inboxes, requests), inboxes, requests
}

func TestGetOrCreateInboxIdempotent(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()

	first, err := svc.GetOrCreateInbox(ctx, "user-1")
	if err != nil {
		t.Fatalf("get-or-create: %v", err)
	}
	second, err := svc.GetOrCreateInbox(ctx, "user-1")
	if err != nil {
		t.Fatalf("get-or-create second: %v", err)
	}
	if first.PublicUUID != second.PublicUUID || first.ID != second.ID {
		t.Fatalf("expected same inbox, got %+v vs %+v", first, second)
	}
	if !strings.HasPrefix(first.PublicPath, "/tools/webhooks/") {
		t.Fatalf("unexpected public path: %s", first.PublicPath)
	}
}

func TestCaptureUnknownUUID(t *testing.T) {
	svc, _, _ := newWebhookService()
	err := svc.Capture(context.Background(), "00000000-0000-4000-8000-000000000000", webhooksvc.CaptureInput{
		Method: "POST",
		URL:    "/tools/webhooks/x",
		Body:   []byte(`{"a":1}`),
	})
	if err == nil || !apperrors.IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestCaptureSizeLimit(t *testing.T) {
	svc, _, requests := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-size")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}

	oversized := bytes.Repeat([]byte("a"), webhooksvc.MaxBodyBytes+1)
	err = svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
		Method:    "POST",
		URL:       "/tools/webhooks/" + inbox.PublicUUID,
		Body:      nil,
		Oversized: true,
	})
	if !webhooksvc.IsBodyTooLarge(err) {
		t.Fatalf("expected body too large, got %v", err)
	}
	if len(requests.ByID) != 1 {
		t.Fatalf("expected error-marker row, got %d", len(requests.ByID))
	}
	for _, req := range requests.ByID {
		if req.CaptureStatus != webhookmodel.CaptureStatusOversized {
			t.Fatalf("expected oversized status, got %s", req.CaptureStatus)
		}
		if len(req.Body) != 0 {
			t.Fatal("oversized capture must not store unbounded body")
		}
	}

	// Direct oversize via body length also rejected.
	err = svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
		Method: "POST",
		URL:    "/tools/webhooks/" + inbox.PublicUUID,
		Body:   oversized,
	})
	if !webhooksvc.IsBodyTooLarge(err) {
		t.Fatalf("expected body too large for large body, got %v", err)
	}
}

func TestReadBodyLimited(t *testing.T) {
	okBody := bytes.Repeat([]byte("b"), webhooksvc.MaxBodyBytes)
	data, oversized, err := webhooksvc.ReadBodyLimited(bytes.NewReader(okBody))
	if err != nil || oversized || len(data) != webhooksvc.MaxBodyBytes {
		t.Fatalf("expected ok read, got len=%d oversized=%v err=%v", len(data), oversized, err)
	}

	big := bytes.Repeat([]byte("c"), webhooksvc.MaxBodyBytes+10)
	data, oversized, err = webhooksvc.ReadBodyLimited(bytes.NewReader(big))
	if err != nil || !oversized || data != nil {
		t.Fatalf("expected oversized, got len=%d oversized=%v err=%v", len(data), oversized, err)
	}
}

func TestCaptureMethodsAndRetention(t *testing.T) {
	svc, inboxes, requests := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-ret")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	for i, method := range methods {
		err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
			Method: method,
			URL:    "/tools/webhooks/" + inbox.PublicUUID + "?i=" + string(rune('0'+i)),
			Body:   []byte(`{"n":` + string(rune('0'+i)) + `}`),
		})
		if err != nil {
			t.Fatalf("capture %s: %v", method, err)
		}
		// Space created_at for deterministic oldest purge.
		time.Sleep(time.Millisecond)
	}
	if len(requests.ByID) != len(methods) {
		t.Fatalf("expected %d requests, got %d", len(methods), len(requests.ByID))
	}

	// Fill to retention cap then one more.
	base := len(requests.ByID)
	for i := base; i < webhooksvc.MaxStoredRequests; i++ {
		if err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
			Method: "POST",
			URL:    "/fill",
			Body:   []byte("x"),
		}); err != nil {
			t.Fatalf("fill capture: %v", err)
		}
		// Avoid identical timestamps for purge order.
		for _, req := range requests.ByID {
			_ = req
		}
	}
	// Force distinct timestamps by rewriting created_at on existing rows.
	n := 0
	for _, req := range requests.ByID {
		req.CreatedAt = time.Unix(1_700_000_000+int64(n), 0).UTC()
		n++
	}

	if err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
		Method: "POST",
		URL:    "/overflow",
		Body:   []byte("overflow"),
	}); err != nil {
		t.Fatalf("overflow capture: %v", err)
	}
	if got := len(requests.ByID); got != webhooksvc.MaxStoredRequests {
		t.Fatalf("expected retention cap %d, got %d", webhooksvc.MaxStoredRequests, got)
	}
	fresh, _ := inboxes.GetByUserID(ctx, "user-ret")
	if fresh.LifetimeReceived != webhooksvc.MaxStoredRequests {
		t.Fatalf("lifetime after purge want %d got %d", webhooksvc.MaxStoredRequests, fresh.LifetimeReceived)
	}
	if fresh.ActiveCount != webhooksvc.MaxStoredRequests {
		t.Fatalf("active after purge want %d got %d", webhooksvc.MaxStoredRequests, fresh.ActiveCount)
	}
}

func TestSoftDeleteCounters(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-del")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
			Method: "POST",
			URL:    "/x",
			Body:   []byte("a"),
		}); err != nil {
			t.Fatalf("capture: %v", err)
		}
	}
	list, err := svc.ListRequests(ctx, "user-del", webhookdto.ListRequestsQuery{Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.ActiveCount != 3 || list.LifetimeReceived != 3 || len(list.Items) != 3 {
		t.Fatalf("unexpected list: %+v", list)
	}
	id := list.Items[0].ID
	del, err := svc.SoftDeleteRequest(ctx, "user-del", id)
	if err != nil {
		t.Fatalf("soft-delete: %v", err)
	}
	if !del.OK || del.ActiveCount != 2 || del.LifetimeReceived != 3 {
		t.Fatalf("unexpected soft-delete response: %+v", del)
	}
	// Idempotent second delete.
	del2, err := svc.SoftDeleteRequest(ctx, "user-del", id)
	if err != nil {
		t.Fatalf("soft-delete again: %v", err)
	}
	if del2.ActiveCount != 2 || del2.LifetimeReceived != 3 {
		t.Fatalf("idempotent soft-delete changed counters: %+v", del2)
	}
	// Missing request ID still returns ok (contract idempotency).
	delMissing, err := svc.SoftDeleteRequest(ctx, "user-del", "00000000-0000-4000-8000-000000000099")
	if err != nil || !delMissing.OK || delMissing.ActiveCount != 2 {
		t.Fatalf("missing soft-delete should be ok: %+v err=%v", delMissing, err)
	}
	list2, err := svc.ListRequests(ctx, "user-del", webhookdto.ListRequestsQuery{Limit: 10})
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(list2.Items) != 2 {
		t.Fatalf("soft-deleted item still listed: %+v", list2.Items)
	}
	_, err = svc.GetRequest(ctx, "user-del", id)
	if err == nil || !apperrors.IsNotFound(err) {
		t.Fatalf("expected not found for soft-deleted detail, got %v", err)
	}
}

func TestSoftDeleteIdempotentWithoutInbox(t *testing.T) {
	svc, _, _ := newWebhookService()
	del, err := svc.SoftDeleteRequest(context.Background(), "user-no-inbox", "00000000-0000-4000-8000-000000000001")
	if err != nil {
		t.Fatalf("expected ok without inbox, got %v", err)
	}
	if !del.OK || del.ActiveCount != 0 || del.LifetimeReceived != 0 {
		t.Fatalf("unexpected response: %+v", del)
	}
}

func TestDetailBodyEncoding(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-body")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}

	cases := []struct {
		name     string
		body     []byte
		encoding string
		binary   bool
	}{
		{name: "empty", body: nil, encoding: "utf-8", binary: false},
		{name: "utf8", body: []byte(`{"a":1}`), encoding: "utf-8", binary: false},
		{name: "binary", body: []byte{0x00, 0xff, 0xfe}, encoding: "base64", binary: true},
	}
	for _, tc := range cases {
		if err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{
			Method: "POST",
			URL:    "/body/" + tc.name,
			Body:   tc.body,
		}); err != nil {
			t.Fatalf("capture %s: %v", tc.name, err)
		}
	}
	list, err := svc.ListRequests(ctx, "user-body", webhookdto.ListRequestsQuery{Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	byURL := map[string]string{}
	for _, item := range list.Items {
		byURL[item.URL] = item.ID
	}
	for _, tc := range cases {
		detail, err := svc.GetRequest(ctx, "user-body", byURL["/body/"+tc.name])
		if err != nil {
			t.Fatalf("detail %s: %v", tc.name, err)
		}
		if detail.BodyEncoding != tc.encoding || detail.IsBinary != tc.binary {
			t.Fatalf("%s: encoding=%q isBinary=%v want %q/%v", tc.name, detail.BodyEncoding, detail.IsBinary, tc.encoding, tc.binary)
		}
		if tc.encoding == "utf-8" && string(tc.body) != detail.Body {
			t.Fatalf("%s: body mismatch %q vs %q", tc.name, detail.Body, tc.body)
		}
		if tc.encoding == "base64" && detail.Body == "" {
			t.Fatalf("%s: expected base64 body", tc.name)
		}
	}
}

func TestRegenerateInvalidatesOldUUID(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-regen")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}
	oldUUID := inbox.PublicUUID
	if err := svc.Capture(ctx, oldUUID, webhooksvc.CaptureInput{Method: "GET", URL: "/old"}); err != nil {
		t.Fatalf("capture before regen: %v", err)
	}
	updated, err := svc.RegenerateUUID(ctx, "user-regen")
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	if updated.PublicUUID == oldUUID {
		t.Fatal("expected new public uuid")
	}
	err = svc.Capture(ctx, oldUUID, webhooksvc.CaptureInput{Method: "GET", URL: "/stale"})
	if err == nil || !apperrors.IsNotFound(err) {
		t.Fatalf("expected old uuid 404, got %v", err)
	}
	if err := svc.Capture(ctx, updated.PublicUUID, webhooksvc.CaptureInput{Method: "POST", URL: "/new", Body: []byte("ok")}); err != nil {
		t.Fatalf("new uuid capture: %v", err)
	}
	list, err := svc.ListRequests(ctx, "user-regen", webhookdto.ListRequestsQuery{Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("history should remain, got %d items", len(list.Items))
	}
}

func TestOwnerAuthorizationNonOwnerCannotList(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "owner")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}
	if err := svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{Method: "POST", URL: "/x", Body: []byte("z")}); err != nil {
		t.Fatalf("capture: %v", err)
	}
	_, err = svc.ListRequests(ctx, "other-user", webhookdto.ListRequestsQuery{})
	if err == nil || !apperrors.IsNotFound(err) {
		t.Fatalf("expected not found for other user without inbox, got %v", err)
	}
	_, err = svc.GetOrCreateInbox(ctx, "other-user")
	if err != nil {
		t.Fatalf("other inbox: %v", err)
	}
	list, err := svc.ListRequests(ctx, "other-user", webhookdto.ListRequestsQuery{})
	if err != nil {
		t.Fatalf("other list: %v", err)
	}
	if len(list.Items) != 0 {
		t.Fatalf("other user must not see owner requests, got %+v", list.Items)
	}
}

func TestListFilterMethodAndQuery(t *testing.T) {
	svc, _, _ := newWebhookService()
	ctx := context.Background()
	inbox, err := svc.GetOrCreateInbox(ctx, "user-filter")
	if err != nil {
		t.Fatalf("inbox: %v", err)
	}
	_ = svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{Method: "GET", URL: "/path/alpha", Body: []byte("hello")})
	_ = svc.Capture(ctx, inbox.PublicUUID, webhooksvc.CaptureInput{Method: "POST", URL: "/path/beta", Body: []byte("world")})

	list, err := svc.ListRequests(ctx, "user-filter", webhookdto.ListRequestsQuery{Method: "POST", Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].Method != "POST" {
		t.Fatalf("method filter failed: %+v", list.Items)
	}
	listQ, err := svc.ListRequests(ctx, "user-filter", webhookdto.ListRequestsQuery{Q: "alpha", Limit: 10})
	if err != nil {
		t.Fatalf("list q: %v", err)
	}
	if len(listQ.Items) != 1 || !strings.Contains(listQ.Items[0].URL, "alpha") {
		t.Fatalf("q filter failed: %+v", listQ.Items)
	}
}
