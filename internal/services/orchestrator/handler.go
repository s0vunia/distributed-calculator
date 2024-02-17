package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

// EndpointExpression обрабатывает и POST, и GET запрос expression
func EndpointExpression(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetExpression(orchestrator)(w, r)
		case http.MethodPost:
			CreateExpression(orchestrator)(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func CreateExpression(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		if err := r.ParseMultipartForm(10 << 20); err != nil { // Set maxMemory to 10 MB
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		expression := r.FormValue("expression")

		// Проверка на наличие выражения
		if expression == "" {
			http.Error(w, "Missing 'expression' field", http.StatusBadRequest)
			return
		}
		idempotencyKey := r.Header.Get("X-Idempotency-Key")
		// Проверка на наличие ключа идемпотентности
		if idempotencyKey == "" {
			http.Error(w, "Missing 'idempotencyKey' field", http.StatusBadRequest)
			return
		}

		if !validateExpression(expression) {
			http.Error(w, "Invalid expression", http.StatusBadRequest)
			return
		}

		// Если уже существует выражение по ключу идемпотентности - возвращаем его
		expressionByKey, err := orchestrator.GetExpressionByKey(ctx, idempotencyKey)
		if expressionByKey != nil {
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				log.Printf("error get expression by key %e", err)
				return
			}
			expressionStruct, _ := json.MarshalIndent(expressionByKey, "", " ")
			fmt.Fprint(w, string(expressionStruct))
			return
		}

		err, idExpression := orchestrator.CreateExpression(ctx, expression, idempotencyKey)
		if err != nil {
			http.Error(w, "error create expression", http.StatusInternalServerError)
			log.Printf("error create expression %e", err)
			return
		}
		// Respond to the client
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, idExpression)
	}
}

func GetExpression(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()

		idPath := strings.TrimPrefix(r.URL.Path, "/expression/")
		// Проверка, что idPath не пуст
		if idPath == "" || idPath == "/" {
			http.Error(w, "Missing ID parameter", http.StatusNotFound)
			return
		}

		// Извлечение ID из пути
		id := path.Base(idPath)

		expression, err := orchestrator.GetExpression(ctx, id)
		if err != nil {
			return
		}
		if expression == nil {
			http.NotFound(w, r)
			return
		}
		expressionStruct, err := json.MarshalIndent(expression, "", " ")
		fmt.Fprint(w, string(expressionStruct))
	}
}
func GetExpressions(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()

		expressions, err := orchestrator.GetExpressions(ctx)
		if err != nil || expressions == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.MarshalIndent(expressions, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(jsonData))
	}
}

func GetAgents(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		agents, _ := orchestrator.GetAgents()
		jsonData, err := json.MarshalIndent(agents, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(jsonData))
	}
}

func GetSubExpressions(orchestrator IOrchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()

		expressions, err := orchestrator.GetSubExpressions(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(expressions)

		jsonData, err := json.MarshalIndent(expressions, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(jsonData))
	}
}
