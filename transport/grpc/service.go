package grpc

import (
	"context"

	"github.com/egiferdians/micro-auth/models"
	"github.com/egiferdians/micro-auth/transport"
	"github.com/egiferdians/micro-auth/util/converter"
	"github.com/go-kit/kit/log"

	"github.com/egiferdians/micro-auth/protobuf/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	login  grpctransport.Handler
	logger log.Logger
}

func NewGRPCServer(serviceEndpoint transport.Endpoints, logger log.Logger) pb.AuthServiceServer {

	var options []grpctransport.ServerOption
	errorLogger := grpctransport.ServerErrorLogger(logger)
	options = append(options, errorLogger)

	return &grpcServer{
		login: grpctransport.NewServer(
			serviceEndpoint.Login,
			decodeLoginRequest,
			encodeLoginResponse,
			options...,
		),
		logger: logger,
	}
}

func (g *grpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginResponse), nil
}

func decodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginRequest)
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func encodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.LoginResponse)
	return &pb.LoginResponse{Status: res.Status, Data: converter.ConvertToPBAuth(res.Data)}, nil
}
