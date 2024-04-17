package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"myproject/internal/models"
	"myproject/internal/repositories"
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

func (p *PostgresRepository) Create(ctx context.Context, login string, passHash []byte) (int64, error) {
	const op = "repositories.user.postgres.Create"
	var lastInsertId int64
	err := p.db.QueryRowContext(ctx, "INSERT INTO users(login, pass_hash) VALUES ($1, $2) RETURNING id", login, string(passHash)).Scan(&lastInsertId)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, repositories.ErrUserExists)
		}
		return 0, fmt.Errorf("create user failure %e", err)
	}
	return lastInsertId, nil
}

func (s *PostgresRepository) Get(ctx context.Context, login string) (models.User, error) {
	const op = "repositories.user.postgres.Get"

	stmt, err := s.db.Prepare("SELECT id, login, pass_hash FROM users WHERE login = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, login)

	var user models.User
	err = row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, repositories.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
