package agent

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"myproject/internal/models"
	"time"
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

func (p *PostgresRepository) Create(s string) error {
	err := p.db.QueryRow("INSERT INTO agents (id) VALUES ($1)", s)
	if err != nil {
		return fmt.Errorf("create agent failure %e", err)
	}
	return nil
}

func (p *PostgresRepository) IsExists(s string) (bool, error) {
	exists, err := rowExists(p.db, "SELECT EXISTS(SELECT 1 FROM agents WHERE id = $1)", s)
	if err != nil {
		return false, fmt.Errorf("agent is exists failure%e", err)
	}
	return exists, nil
}

func (p *PostgresRepository) UpdateHeartbeat(id string) error {
	now := time.Now()
	_, err := p.db.Exec("UPDATE agents SET heartbeat=$1 WHERE id=$2",
		now, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) CreateIfNotExistsAndUpdateHeartbeat(id string) error {
	isExists, _ := p.IsExists(id)
	if !isExists {
		_ = p.Create(id)
	} else {
		_ = p.UpdateHeartbeat(id)
	}
	return nil
}

func (p *PostgresRepository) GetAgents() ([]*models.Agent, error) {
	rows, err := p.db.Query("SELECT id, heartbeat FROM agents")
	if err != nil {
		log.Printf("error get agents query")
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		var timestamp time.Time
		var agent models.Agent
		if err := rows.Scan(&agent.Id, &timestamp); err != nil {
			log.Printf("error scan agent")
			return nil, err
		}
		agent.Heartbeat = timestamp.Unix()
		agents = append(agents, &agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

// rowExists универсальная функция для проверки записей на существование
func rowExists(db *sql.DB, query string, args ...interface{}) (bool, error) {
	var exists bool
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("error checking if row exists: %v", err)
	}
	return exists, nil
}
