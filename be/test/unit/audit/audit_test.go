package audit_test

import (
	"context"
	"testing"
	"time"

	"be/internal/audit"
)

func TestBuildUsesContextActorAndRequest(t *testing.T) {
	ctx := audit.WithRequest(context.Background(), "203.0.113.5", "Agent/1", "req-1")
	ctx = audit.WithActor(ctx, "u-1", "admin@example.com")

	row := audit.Build(ctx, audit.Entry{Action: "user.update", ResourceType: "user", ResourceID: "u-2"}, time.Unix(0, 0))
	if row.ActorID == nil || *row.ActorID != "u-1" || row.ActorEmail != "admin@example.com" ||
		row.IPAddress != "203.0.113.5" || row.UserAgent != "Agent/1" || row.RequestID != "req-1" || row.Metadata == nil {
		t.Fatalf("unexpected row: %+v", row)
	}

	// Explicit actor (sign-in flows) wins over the context.
	row = audit.Build(ctx, audit.Entry{Action: "auth.login", ActorID: "u-9", ActorEmail: "x@example.com"}, time.Now())
	if *row.ActorID != "u-9" || row.ActorEmail != "x@example.com" {
		t.Fatalf("override ignored: %+v", row)
	}

	row = audit.Build(context.Background(), audit.Entry{Action: "system"}, time.Now())
	if row.ActorID != nil {
		t.Fatal("anonymous entries must have a nil actor id (FK)")
	}
}
