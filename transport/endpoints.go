package transport

import (
	"context"
	"fmt"

	"github.com/egiferdians/micro-auth/models"
	svc "github.com/egiferdians/micro-auth/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login        endpoint.Endpoint
}

func MakeEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		Login:        makeLoginEndpoint(s),
	}
}

func makeLoginEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*models.User)
		status,_, authentic, err := s.Login(ctx, req.Email, req.Password)
		fmt.Println("authentic ", authentic)
		return LoginResponse{Status: status, Data: authentic}, err
	}
}