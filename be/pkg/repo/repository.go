package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"be/pkg/postgres"
	"be/pkg/query"
)

// Opts configures a table-scoped generic repository.
type Opts struct {
	Table            string
	PK               string
	SoftDeleteColumn string
}

// Repository is a reusable CRUD helper over pgx + squirrel.
type Repository[T any] struct {
	pg   *postgres.Postgres
	opts Opts
}

// New constructs a repository for table T.
func New[T any](pg *postgres.Postgres, opts Opts) *Repository[T] {
	if strings.TrimSpace(opts.PK) == "" {
		opts.PK = "id"
	}
	return &Repository[T]{pg: pg, opts: opts}
}

// DB returns the shared Postgres client.
func (r *Repository[T]) DB() *postgres.Postgres {
	return r.pg
}

func (r *Repository[T]) table() string {
	return postgres.QuoteIdent(r.opts.Table)
}

func (r *Repository[T]) pk() string {
	return postgres.QuoteIdent(r.opts.PK)
}

func (r *Repository[T]) applySoftDelete(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
	if r.opts.SoftDeleteColumn == "" {
		return sb
	}
	return sb.Where(squirrel.Expr(postgres.QuoteIdent(r.opts.SoftDeleteColumn) + " IS NULL"))
}

// Select starts SELECT * FROM table (plus default soft-delete filter).
func (r *Repository[T]) Select() squirrel.SelectBuilder {
	cols := columnSelectExprs[T]()
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	return r.applySoftDelete(r.pg.Builder.Select(cols...).From(r.table()))
}

// FindByID loads one row by primary key. Missing rows return (nil, nil).
func (r *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}
	sb := r.Select().Where(squirrel.Expr(r.pk()+" = ?", id)).Limit(1)
	return r.QueryOne(ctx, sb)
}

// FindOne loads the first row matching q. Missing rows return (nil, nil).
func (r *Repository[T]) FindOne(ctx context.Context, q *query.Query) (*T, error) {
	sb := query.ApplyOrderLimit(query.ApplySelect(r.Select(), q), q)
	sb = sb.Limit(1)
	return r.QueryOne(ctx, sb)
}

// Find loads all rows matching q (pagination from q is applied).
func (r *Repository[T]) Find(ctx context.Context, q *query.Query) ([]T, error) {
	sb := query.ApplyOrderLimit(query.ApplySelect(r.Select(), q), q)
	return r.Query(ctx, sb)
}

// Paginate counts then loads a page. Returns items and total matching rows.
func (r *Repository[T]) Paginate(ctx context.Context, q *query.Query) ([]T, int64, error) {
	if q == nil {
		q = query.New(query.DefaultPage, query.DefaultPageSize)
	}

	countSB := query.ApplySelect(r.applySoftDelete(r.pg.Builder.Select("COUNT(*)").From(r.table())), q)
	total, err := r.QueryCount(ctx, countSB)
	if err != nil {
		return nil, 0, err
	}

	sb := query.ApplyOrderLimit(query.ApplySelect(r.Select(), q), q)
	items, err := r.Query(ctx, sb)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Insert writes entity and scans RETURNING * back into it (fills generated IDs).
func (r *Repository[T]) Insert(ctx context.Context, entity *T) error {
	if entity == nil {
		return fmt.Errorf("insert: nil entity")
	}
	cols, err := columnMap(entity, columnMapOpts{
		PK:           r.opts.PK,
		SkipEmptyPK:  true,
		SkipZeroTime: true,
		SkipNilPtr:   true,
	})
	if err != nil {
		return err
	}
	if len(cols) == 0 {
		return fmt.Errorf("insert %s: no columns", r.opts.Table)
	}

	sql, args, err := r.pg.Builder.Insert(r.table()).SetMap(cols).Suffix("RETURNING " + r.pk()).ToSql()
	if err != nil {
		return fmt.Errorf("insert build: %w", err)
	}
	var id string
	if err := r.pg.Querier(ctx).QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return err
	}
	if strings.TrimSpace(id) != "" {
		setStringColumn(entity, r.opts.PK, id)
	}
	return nil
}

// Update writes all mapped columns for entity keyed by PK.
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return fmt.Errorf("update: nil entity")
	}
	pkValue, cols, err := updateMap(entity, r.opts.PK)
	if err != nil {
		return err
	}
	if pkValue == "" {
		return fmt.Errorf("update %s: empty primary key", r.opts.Table)
	}
	sql, args, err := r.pg.Builder.Update(r.table()).
		SetMap(cols).
		Where(squirrel.Expr(r.pk()+" = ?", pkValue)).
		ToSql()
	if err != nil {
		return fmt.Errorf("update build: %w", err)
	}
	_, err = r.pg.Querier(ctx).Exec(ctx, sql, args...)
	return err
}

// UpdateMap patches columns for one primary key.
func (r *Repository[T]) UpdateMap(ctx context.Context, id string, cols map[string]any) (int64, error) {
	if strings.TrimSpace(id) == "" || len(cols) == 0 {
		return 0, nil
	}
	quoted := make(map[string]any, len(cols))
	for key, value := range cols {
		quoted[postgres.QuoteIdent(key)] = value
	}
	sql, args, err := r.pg.Builder.Update(r.table()).
		SetMap(quoted).
		Where(squirrel.Expr(r.pk()+" = ?", id)).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("update map build: %w", err)
	}
	tag, err := r.pg.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// UpdateWhere patches columns matching q. Returns rows affected.
func (r *Repository[T]) UpdateWhere(ctx context.Context, q *query.Query, cols map[string]any) (int64, error) {
	if len(cols) == 0 {
		return 0, nil
	}
	quoted := make(map[string]any, len(cols))
	for key, value := range cols {
		quoted[postgres.QuoteIdent(key)] = value
	}
	ub := r.pg.Builder.Update(r.table()).SetMap(quoted)
	if sqlizer := query.GroupSqlizer(q.Conditions); sqlizer != nil {
		ub = ub.Where(sqlizer)
	}
	sql, args, err := ub.ToSql()
	if err != nil {
		return 0, fmt.Errorf("update where build: %w", err)
	}
	tag, err := r.pg.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// DeleteByID hard-deletes, or sets SoftDeleteColumn when configured.
func (r *Repository[T]) DeleteByID(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return nil
	}
	if r.opts.SoftDeleteColumn != "" {
		_, err := r.UpdateMap(ctx, id, map[string]any{r.opts.SoftDeleteColumn: nowUTC()})
		return err
	}
	sql, args, err := r.pg.Builder.Delete(r.table()).Where(squirrel.Expr(r.pk()+" = ?", id)).ToSql()
	if err != nil {
		return fmt.Errorf("delete build: %w", err)
	}
	_, err = r.pg.Querier(ctx).Exec(ctx, sql, args...)
	return err
}

// Count returns matching row count (soft-delete applied).
func (r *Repository[T]) Count(ctx context.Context, q *query.Query) (int64, error) {
	sb := query.ApplySelect(r.applySoftDelete(r.pg.Builder.Select("COUNT(*)").From(r.table())), q)
	return r.QueryCount(ctx, sb)
}

// Query executes a select builder into []T.
func (r *Repository[T]) Query(ctx context.Context, sb squirrel.SelectBuilder) ([]T, error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query build: %w", err)
	}
	rows, err := r.pg.Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []T{}
	}
	return items, nil
}

// QueryOne executes a select builder into one T. Missing rows return (nil, nil).
func (r *Repository[T]) QueryOne(ctx context.Context, sb squirrel.SelectBuilder) (*T, error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query one build: %w", err)
	}
	var dest T
	if err := r.scanOneInto(ctx, sql, args, &dest); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dest, nil
}

// QueryCount executes a COUNT(*) builder.
func (r *Repository[T]) QueryCount(ctx context.Context, sb squirrel.SelectBuilder) (int64, error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return 0, fmt.Errorf("count build: %w", err)
	}
	var total int64
	if err := r.pg.Querier(ctx).QueryRow(ctx, sql, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

// ExecBuilder runs an insert/update/delete builder and returns rows affected.
func (r *Repository[T]) ExecBuilder(ctx context.Context, b squirrel.Sqlizer) (int64, error) {
	sql, args, err := b.ToSql()
	if err != nil {
		return 0, fmt.Errorf("exec build: %w", err)
	}
	tag, err := r.pg.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (r *Repository[T]) scanOneInto(ctx context.Context, sql string, args []any, dest *T) error {
	rows, err := r.pg.Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[T])
	if err != nil {
		return err
	}
	*dest = item
	return nil
}
