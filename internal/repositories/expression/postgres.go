package expression

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"myproject/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dataSourceName string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) CreateExpression(ctx context.Context, s, idempotencyId string) (*models.Expression, error) {
	var id string
	expression := &models.Expression{
		IdempotencyKey: idempotencyId,
		Value:          s,
		State:          models.ExpressionState(models.InProgress),
	}

	err := r.db.QueryRowContext(ctx, "INSERT INTO expressions (id, idempotency_key, value, state) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id",
		expression.IdempotencyKey, expression.Value, expression.State).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("create expression failure %e", err)
	}

	expression.Id = id
	return expression, nil
}

func (r *PostgresRepository) GetExpressions(ctx context.Context) ([]*models.Expression, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, idempotency_key, value, state, result FROM expressions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []*models.Expression
	for rows.Next() {
		var expr models.Expression
		var result sql.NullFloat64
		if err := rows.Scan(&expr.Id, &expr.IdempotencyKey, &expr.Value, &expr.State, &result); err != nil {
			return nil, err
		}
		if result.Valid {
			expr.Result = result.Float64
		} else {
			expr.Result = 0 // или любое другое значение по умолчанию
		}
		expressions = append(expressions, &expr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expressions, nil
}

func (r *PostgresRepository) GetExpressionById(ctx context.Context, id string) (*models.Expression, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, idempotency_key, value, state, result FROM expressions WHERE id=$1", id)
	var expr models.Expression
	var result sql.NullFloat64
	if err := row.Scan(&expr.Id, &expr.IdempotencyKey, &expr.Value, &expr.State, &result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if result.Valid {
		expr.Result = result.Float64
	} else {
		expr.Result = 0 // или любое другое значение по умолчанию
	}

	return &expr, nil
}

func (r *PostgresRepository) GetExpressionByKey(ctx context.Context, key string) (*models.Expression, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, idempotency_key, value, state, result FROM expressions WHERE idempotency_key=$1", key)
	var expr models.Expression
	var result sql.NullFloat64
	if err := row.Scan(&expr.Id, &expr.IdempotencyKey, &expr.Value, &expr.State, &result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if result.Valid {
		expr.Result = result.Float64
	} else {
		expr.Result = 0 // или любое другое значение по умолчанию
	}

	return &expr, nil
}

func (r *PostgresRepository) UpdateExpression(ctx context.Context, expression *models.Expression) error {
	_, err := r.db.ExecContext(ctx, "UPDATE expressions SET state=$1, result=$2 WHERE id=$3",
		expression.State, expression.Result, expression.Id)
	return err
}

func (r *PostgresRepository) UpdateExpressionById(ctx context.Context, id uuid.UUID, result float64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE expressions SET result=$1, state=$3 WHERE id=$2",
		result, id, models.ExpressionState(models.Ok))
	return err
}

func (r *PostgresRepository) UpdateState(ctx context.Context, key string, state models.ExpressionState) error {
	_, err := r.db.ExecContext(ctx, "UPDATE expressions SET state=$2 WHERE id=$1",
		key, state)
	return err
}

func (r *PostgresRepository) DeleteExpressionById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM expressions WHERE id=$1",
		id.String())
	if err != nil {
		return err
	}
	return nil
}

// Close closes the database connection.
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
