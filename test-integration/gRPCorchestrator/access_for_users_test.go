package gRPCorchestrator

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	protos "github.com/s0vunia/protos/gen/go/auth"
	"github.com/s0vunia/protos/gen/go/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"testing"
)

func TestGRPCServiceAccessForDiffUsers(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//cfg := config.MustLoadPath("../../config/local_tests.yaml")

	conn, err := grpc.Dial("localhost:44044", grpc.WithInsecure())
	defer conn.Close()
	assert.NoError(t, err)

	// authenticated
	logins, passwords := []string{"test1", "test2"}, []string{"testpass1", "testpass2"}
	tokens := make([]string, 0, 2)
	for i := 0; i < len(tokens); i++ {
		authClient := protos.NewAuthClient(conn)
		registerResponse, err := authClient.Register(context.Background(), &protos.RegisterRequest{
			Login:    logins[i],
			Password: passwords[i],
		})
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.OK {
			assert.Equal(t, codes.AlreadyExists, st.Code())
		}
		log.Printf("%v", registerResponse)

		loginResponse, err := authClient.Login(context.Background(), &protos.LoginRequest{
			Login:    logins[i],
			Password: passwords[i],
			AppId:    1,
		})
		assert.NoError(t, err)
		token := loginResponse.Token
		tokens = append(tokens, token)
	}

	client := orchestrator.NewOrchestratorClient(conn)

	var expressionsIds []string
	for _, token := range tokens {
		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		ctx = context.WithValue(ctx, "userID", 1)
		createExpressionResponse, err := client.CreateExpression(ctx, &orchestrator.CreateExpressionRequest{
			IdempotencyKey: uuid.New().String(),
			Expression:     "2+2*2",
		})
		assert.NoError(t, err)
		log.Printf("%v", createExpressionResponse)
		expressionsIds = append(expressionsIds, createExpressionResponse.ExpressionId)
	}

	for i, tokenString := range tokens {
		for j, expressionId := range expressionsIds {
			md := metadata.New(map[string]string{
				"authorization": tokenString,
			})
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			assert.NoError(t, err)
			claims, ok := token.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			userId, ok := claims["uid"].(float64)
			assert.True(t, ok)
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			ctx = context.WithValue(ctx, "userID", userId)

			_, err = client.GetExpression(ctx, &orchestrator.GetExpressionRequest{
				ExpressionId: expressionId,
			})
			if i == j {
				log.Println(err)
				assert.NoError(t, err)
			} else {
				log.Println(err)
				assert.Error(t, err)
			}
		}
	}
}
