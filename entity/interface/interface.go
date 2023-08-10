package _interface

import (
	"context"

	"github.com/egiferdians/micro-auth/auth"
	"github.com/egiferdians/micro-auth/models"
)

type Repository interface {
	Login(ctx context.Context, usr models.User) (*models.User, *auth.Authenticated, error)
}