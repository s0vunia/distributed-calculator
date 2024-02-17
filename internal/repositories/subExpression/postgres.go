package subExpression

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"myproject/project/internal/models"
	"time"
)

type PostgresRepository struct {
	db                     *sql.DB
	listenerSubExpressions *pq.Listener
	listener               chan *models.SubExpression
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

	listener := pq.NewListener(dataSourceName, 1*time.Second, 5*time.Second, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err)
		}
	})

	// Подписываемся на канал уведомлений
	err = listener.Listen("sub_expressions_channel")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init listener: %w", err)
	}

	// Запускаем горутину для обработки уведомлений
	listenerModel := make(chan *models.SubExpression)

	repo := &PostgresRepository{db, listener, listenerModel}
	//go repo.handleNotificationNotTrigger(context.Background(), listenerModel)
	go repo.handleNotification(context.Background(), listener.Notify, listenerModel)
	return repo, nil
}

func (r *PostgresRepository) CreateSubExpression(ctx context.Context, subExpression *models.SubExpression) (*models.SubExpression, error) {
	var id uuid.UUID

	err := r.db.QueryRowContext(ctx, "INSERT INTO sub_expressions (id, expressions_id, val1, val2,sub_expression_id1,sub_expression_id2,is_last, action, error) VALUES (gen_random_uuid(), $1, $2, $3, NULLIF($4, '')::UUID, NULLIF($5, '')::UUID, $6, $7, $8) RETURNING id",
		subExpression.ExpressionId, subExpression.Val1, subExpression.Val2, subExpression.SubExpressionId1, subExpression.SubExpressionId2, subExpression.IsLast, subExpression.Action, subExpression.Error).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("create expression failure %e", err)
	}

	subExpression.Id = id
	return subExpression, nil
}

// handleNotification обрабатывает триггеры бд (триггер по обновлению и добавлению subexpressions)
func (r *PostgresRepository) handleNotification(ctx context.Context, notificationChan <-chan *pq.Notification,
	listener chan *models.SubExpression) {
	defer close(listener)
	for {
		select {
		case notification := <-notificationChan:
			if notification == nil || notification.Extra == "" {
				continue
			}
			var se *models.SubExpression
			err := json.Unmarshal([]byte(notification.Extra), &se)
			if err != nil {
				log.Printf("Failed to parse JSON: %v", err)
				continue
			}
			newse, err := r.GetExpressionByKey(ctx, se.Id.String())
			if err != nil {
				log.Printf("Failed to get expr by key: %s %v %e", se.Id.String(), se, err)
				continue
			}
			listener <- newse
		}
	}
}

func (r *PostgresRepository) GetSubExpressions() chan *models.SubExpression {
	return r.listener
}

func (r *PostgresRepository) GetSubExpressionsList(ctx context.Context) ([]*models.SubExpression, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, expressions_id, sub_expression_id1, sub_expression_id2 FROM sub_expressions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []*models.SubExpression
	for rows.Next() {
		var expr models.SubExpression
		var result sql.NullFloat64
		if err := rows.Scan(&expr.Id, &expr.ExpressionId, &expr.SubExpressionId1, &expr.SubExpressionId2); err != nil {
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

func (r *PostgresRepository) UpdateSubExpressions(ctx context.Context, expression *models.SubExpression) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sub_expressions SET result=$1 WHERE id=$2",
		expression.Result, expression.Id)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, "UPDATE sub_expressions\nSET val1 = CASE WHEN sub_expression_id1 = $1 THEN $2 ELSE val1 END, val2 = CASE WHEN sub_expression_id2 = $1 THEN $2 ELSE val2 END\nWHERE sub_expression_id1 = $1 OR sub_expression_id2 = $1;\n",
		expression.Id, expression.Result)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, "UPDATE sub_expressions SET sub_expression_id1 = NULL WHERE sub_expression_id1 = $1",
		expression.Id)
	_, err = r.db.ExecContext(ctx, "UPDATE sub_expressions SET sub_expression_id2 = NULL WHERE sub_expression_id2 = $1",
		expression.Id)
	return err
}

func (r *PostgresRepository) GetExpressionByKey(ctx context.Context, key string) (*models.SubExpression, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT id, expressions_id, val1, val2, sub_expression_id1, sub_expression_id2, action, is_last, error FROM sub_expressions WHERE id = $1",
		key)
	var expr models.SubExpression
	var result sql.NullFloat64
	if err := rows.Scan(&expr.Id, &expr.ExpressionId, &expr.Val1, &expr.Val2, &expr.SubExpressionId1, &expr.SubExpressionId2, &expr.Action, &expr.IsLast, &expr.Error); err != nil {
		return nil, err
	}
	if result.Valid {
		expr.Result = result.Float64
	} else {
		expr.Result = 0 // или любое другое значение по умолчанию
	}
	return &expr, nil
}

func (r *PostgresRepository) DeleteSubExpressionsByExpressionId(ctx context.Context, expressionId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sub_expressions WHERE expressions_id=$1",
		expressionId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) UpdateSubExpressionAgent(ctx context.Context, idSubExpression, agentId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sub_expressions SET agent_id=$1 WHERE id=$2",
		agentId, idSubExpression)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) DeleteSubExpressionById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sub_expressions WHERE id=$1",
		id.String())
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetNotCalculatedSubExpressionsByAgentId(ctx context.Context, agentId uuid.UUID) ([]*models.SubExpression, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, expressions_id, sub_expression_id1, sub_expression_id2, val1, val2, action, is_last, error FROM sub_expressions WHERE agent_id=$1 AND result IS NULL",
		agentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []*models.SubExpression
	for rows.Next() {
		var expr models.SubExpression
		var result sql.NullFloat64
		if err := rows.Scan(&expr.Id, &expr.ExpressionId, &expr.SubExpressionId1, &expr.SubExpressionId2, &expr.Val1, &expr.Val2, &expr.Action, &expr.IsLast, &expr.Error); err != nil {
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

func (r *PostgresRepository) ReplaceExpressionsIds(ctx context.Context, oldId uuid.UUID, newId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sub_expressions SET sub_expression_id1=$1 WHERE sub_expression_id1=$2",
		newId.String(), oldId.String())
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, "UPDATE sub_expressions SET sub_expression_id2=$1 WHERE sub_expression_id2=$2",
		newId.String(), oldId.String())
	if err != nil {
		return err
	}
	return nil
}
