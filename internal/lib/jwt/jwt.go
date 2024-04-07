package jwt

import (
	"context"
	"fmt"
	"myproject/internal/models"
	"myproject/internal/repositories/app"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken creates new JWT token for given user and app.
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ProcessJWT Функция для извлечения app_id из JWT и проверки его валидности
func ProcessJWT(ctx context.Context, tokenString string, appRepo app.Repository) error {
	// Извлечение app_id из JWT без проверки подписи
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid JWT")
	}

	appID, ok := claims["app_id"].(string)
	if !ok {
		return fmt.Errorf("app_id not found in JWT")
	}
	appIdString, err := strconv.Atoi(appID)
	app, err := appRepo.App(ctx, appIdString)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	// Проверка валидности JWT с использованием полученного секрета
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.Secret), nil
	})

	if err != nil {
		return fmt.Errorf("invalid JWT: %w", err)
	}
	return nil
}
