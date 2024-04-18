package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"myproject/internal/models"
	"myproject/internal/repositories/app"
	mocksApp "myproject/internal/repositories/app/mocks"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	type args struct {
		user     models.User
		app      models.App
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				user: models.User{
					ID:       1234,
					Login:    "krutoy",
					PassHash: []byte("hahaha"),
				},
				app: models.App{
					Secret: "secret",
					Name:   "приложуха",
					ID:     1,
				},
				duration: time.Second * 10,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				user: models.User{
					ID:       1234,
					Login:    "hoe",
					PassHash: []byte("hahahahi"),
				},
				app: models.App{
					Secret: "",
					Name:   "приложуха",
					ID:     1,
				},
				duration: time.Second * 10,
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				user: models.User{
					ID:       1234,
					Login:    "hoe",
					PassHash: []byte("hahahahi"),
				},
				app: models.App{
					Secret: "",
					Name:   "приложуха",
					ID:     1,
				},
				duration: -time.Second * 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewToken(tt.args.user, tt.args.app, tt.args.duration)
			claims := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(got, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.args.app.Secret), nil
			})
			if !tt.wantErr {
				assert.NoError(t, err)
				appID, ok := claims["app_id"].(float64)
				assert.True(t, ok)
				assert.Equal(t, appID, float64(tt.args.app.ID))

				uid, ok := claims["uid"].(float64)
				assert.True(t, ok)
				assert.Equal(t, uid, float64(tt.args.user.ID))

				login, ok := claims["login"].(string)
				assert.True(t, ok)
				assert.Equal(t, login, tt.args.user.Login)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestProcessJWT(t *testing.T) {
	type args struct {
		ctx         context.Context
		tokenString string
		appRepo     app.Repository
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := models.App{
				ID:     1,
				Name:   "testing",
				Secret: "testing_secret",
			}
			token, err := NewToken(
				models.User{
					ID:       1,
					Login:    "testing_login",
					PassHash: []byte("testhash"),
				},
				app,
				time.Second*10,
			)
			assert.NoError(t, err)
			tt.args.tokenString = token

			appRepo := mocksApp.NewRepository(t)
			appRepo.
				On("App", mock.Anything, app.ID).
				Return(app, nil)

			err = ProcessJWT(tt.args.ctx, tt.args.tokenString, appRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
