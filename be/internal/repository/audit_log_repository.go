package repository

import (
	"context"
	"time"

	auditmodel "be/internal/models/audit"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type AuditLogRepository struct {
	*repo.Repository[auditmodel.Log]
	db *postgres.Postgres
}

var _ interfaces.AuditLogRepository = (*AuditLogRepository)(nil)

func NewAuditLogRepository(db *postgres.Postgres) *AuditLogRepository {
	return &AuditLogRepository{
		Repository: repo.New[auditmodel.Log](db, repo.Opts{Table: "audit_logs", PK: "id"}),
		db:         db,
	}
}

func (r *AuditLogRepository) Insert(ctx context.Context, log *auditmodel.Log) error {
	return r.Repository.Insert(ctx, log)
}

func (r *AuditLogRepository) List(ctx context.Context, q *query.Query) ([]auditmodel.Log, int64, error) {
	return r.Paginate(ctx, q)
}

func (r *AuditLogRepository) CountByActionPerDay(ctx context.Context, action string, since time.Time) (map[string]int64, error) {
	rows, err := r.db.Querier(ctx).Query(ctx, `
SELECT to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD') AS day, COUNT(*)
FROM audit_logs WHERE action = $1 AND created_at >= $2
GROUP BY day`, action, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var day string
		var n int64
		if err := rows.Scan(&day, &n); err != nil {
			return nil, err
		}
		out[day] = n
	}
	return out, rows.Err()
}

func (r *AuditLogRepository) CountByAction(ctx context.Context, action string, from, to time.Time) (int64, error) {
	var n int64
	err := r.db.Querier(ctx).QueryRow(ctx,
		`SELECT COUNT(*) FROM audit_logs WHERE action = $1 AND created_at >= $2 AND created_at < $3`,
		action, from, to).Scan(&n)
	return n, err
}

func (r *AuditLogRepository) DeleteBefore(ctx context.Context, t time.Time) (int64, error) {
	tag, err := r.db.Querier(ctx).Exec(ctx, `DELETE FROM audit_logs WHERE created_at < $1`, t)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
