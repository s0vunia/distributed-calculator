package authgrpc

import (
	"context"
	gojwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"myproject/internal/lib/jwt"
	"myproject/internal/repositories/app"
)

func JWTMiddleware(appRepo app.Repository) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := values[0]

		// Проверка токена
		// Здесь должен быть ваш код для проверки токена
		// Например, вы можете использовать библиотеку для работы с JWT
		// Если токен недействителен, верните ошибку
		err, jwtToken := jwt.ProcessJWT(ctx, token, appRepo)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		if claims, ok := jwtToken.Claims.(gojwt.MapClaims); ok {
			userId, ok := claims["uid"].(float64)
			if ok {
				// Извлечение данных из токена
				ctx = context.WithValue(ctx, "userID", userId)
			}
		}
		// Если токен действителен, продолжайте обработку запроса
		return handler(ctx, req)
	}
}
