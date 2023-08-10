package usecase

import (
	"context"
	"errors"

	"github.com/egiferdians/micro-auth/auth"
	_interface "github.com/egiferdians/micro-auth/entity/interface"
	"github.com/egiferdians/micro-auth/models"
	svc "github.com/egiferdians/micro-auth/service"
	"github.com/egiferdians/micro-auth/util/errcode"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
)

type microauthService struct {
	repository _interface.Repository
}

func NewMicroAuthService(repo _interface.Repository) svc.Service {
	return &microauthService{
		repository: repo,
	}
}

func (au microauthService) Login(ctx context.Context, email, password string) (string, *models.User, *auth.Authenticated, error) {
	dataNode := models.User{
		Email:        email,
		Password: password,
	}

	err := dataNode.Validate()
	if err != nil {
		return "Failed",nil, nil,  errcode.New(codes.InvalidArgument, err)
	}
	
	uSr, authenticated, err := au.repository.Login(ctx, dataNode)
	if err != nil {
		return "Failed",nil, nil, err
	}
	err = auth.VerifyPassword(uSr.Password, password)
	if err != nil {
		return "Failed",nil, nil, errors.New("Email or Password is incorrect")
	}
	token, refresh, err := auth.GenerateToken(uuid.UUID(uSr.IDUser))
	if err != nil {
		return "Failed",nil, nil, err
	}
	authenticated.User = uSr
	authenticated.AccessToken = token
	authenticated.RefreshToken = refresh
	return "Success", uSr, authenticated, nil
}