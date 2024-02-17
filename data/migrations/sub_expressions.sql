CREATE TABLE IF NOT EXISTS sub_expressions
(
    id                 UUID PRIMARY KEY,
    expressions_id     UUID,
    val1               DOUBLE PRECISION,
    val2               DOUBLE PRECISION,
    sub_expression_id1 UUID,
    sub_expression_id2 UUID,
    action             VARCHAR(50),
    result             DOUBLE PRECISION,
    is_last            BOOL,
    error              BOOL,
    agent_id           UUID,
    created_at timestamp NOT NULL DEFAULT NOW()
);

-- Функция для отправки уведомлений
CREATE OR REPLACE FUNCTION notify_sub_expression_fields()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.sub_expression_id1 IS NULL AND NEW.sub_expression_id2 IS NULL AND NEW.result IS NULL THEN
       PERFORM pg_notify('sub_expressions_channel', json_build_object(
            'id', NEW.id::text,
            'expressions_id', NEW.expressions_id::text,
            'val1', NEW.val1,
            'val2', NEW.val2,
            'sub_expression_id1', coalesce(NEW.sub_expression_id1::text, 'NULL'),
            'sub_expression_id2', coalesce(NEW.sub_expression_id2::text, 'NULL'),
            'action', NEW.action,
            'is_last', NEW.action,
            'result', NEW.result
    )::text);
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER sub_expression_trigger_update
AFTER UPDATE ON sub_expressions
FOR EACH ROW
WHEN (OLD.val1 IS DISTINCT FROM NEW.val1 OR
      OLD.val2 IS DISTINCT FROM NEW.val2 OR
      OLD.sub_expression_id1 IS DISTINCT FROM NEW.sub_expression_id1 OR
      OLD.sub_expression_id2 IS DISTINCT FROM NEW.sub_expression_id2)
EXECUTE PROCEDURE notify_sub_expression_fields();

CREATE TRIGGER sub_expression_trigger_insert
AFTER INSERT ON sub_expressions
FOR EACH ROW EXECUTE PROCEDURE notify_sub_expression_fields()