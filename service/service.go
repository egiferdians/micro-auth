package service

import (
	"context"

	"github.com/egiferdians/micro-auth/auth"
	"github.com/egiferdians/micro-auth/models"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, *models.User, *auth.Authenticated, error)
}