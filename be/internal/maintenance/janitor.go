// Package maintenance runs periodic housekeeping in the API process: pruning
// old audit entries and dead tokens. Every task is idempotent, so running it on
// several replicas at once is harmless.
package maintenance

import (
	"context"
	"log/slog"
	"time"

	"be/internal/settings"
	"be/pkg/postgres"
)

// AuditPruner deletes audit entries older than a cutoff.
type AuditPruner interface {
	DeleteBefore(ctx context.Context, t time.Time) (int64, error)
}

type Janitor struct {
	db       *postgres.Postgres
	audit    AuditPruner
	settings *settings.Service
	interval time.Duration
	logger   *slog.Logger
}

func NewJanitor(db *postgres.Postgres, audit AuditPruner, settingsSvc *settings.Service) *Janitor {
	return &Janitor{db: db, audit: audit, settings: settingsSvc, interval: 6 * time.Hour, logger: slog.Default()}
}

// Run cleans once immediately, then every interval until ctx is done.
func (j *Janitor) Run(ctx context.Context) {
	j.RunOnce(ctx)
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.RunOnce(ctx)
		}
	}
}

// RunOnce performs one cleanup pass and logs what it removed.
func (j *Janitor) RunOnce(ctx context.Context) {
	now := time.Now()
	days := j.settings.Int(ctx, settings.KeyAuditRetentionDays)
	if days > 0 {
		if n, err := j.audit.DeleteBefore(ctx, now.AddDate(0, 0, -days)); err != nil {
			j.logger.Error("janitor: prune audit logs", "err", err)
		} else if n > 0 {
			j.logger.Info("janitor: pruned audit logs", "deleted", n, "retention_days", days)
		}
	}
	j.exec(ctx, "user tokens",
		`DELETE FROM user_tokens WHERE (used_at IS NOT NULL OR expires_at < NOW()) AND created_at < $1`,
		now.AddDate(0, 0, -7))
	j.exec(ctx, "refresh tokens",
		`DELETE FROM refresh_tokens WHERE (revoked_at IS NOT NULL OR expires_at < NOW()) AND updated_at < $1`,
		now.AddDate(0, 0, -30))
}

func (j *Janitor) exec(ctx context.Context, what, sql string, cutoff time.Time) {
	tag, err := j.db.Querier(ctx).Exec(ctx, sql, cutoff)
	if err != nil {
		j.logger.Error("janitor: prune "+what, "err", err)
		return
	}
	if n := tag.RowsAffected(); n > 0 {
		j.logger.Info("janitor: pruned "+what, "deleted", n)
	}
}
