package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/afex/hystrix-go/hystrix"
	mysql "github.com/egiferdians/micro-auth/entity/mysql"
	"github.com/egiferdians/micro-auth/middleware"
	"github.com/egiferdians/micro-auth/protobuf/pb"
	svc "github.com/egiferdians/micro-auth/service"
	"github.com/egiferdians/micro-auth/transport"
	transportgrpc "github.com/egiferdians/micro-auth/transport/grpc"
	usecase "github.com/egiferdians/micro-auth/usecase"
	"github.com/egiferdians/micro/util/dbconnector"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var logger log.Logger
var microAuthHystrixCommand = "blog_command"
var serviceName = "blog"

// logs keys
var (
	MessageLog             = "message"
	DatabaseConnectFailLog = "error_database_connection"
	DatabasePingFailLog    = "error_pinging_database"
	LoadEnvFailLog         = "error_load_env"
)

func loadEnv() {
	// load .env content file
	err := godotenv.Load()
	if err != nil {
		_ = level.Error(logger).Log(LoadEnvFailLog, err)
		os.Exit(-1)
	}
}

// createLogger this method used for initialize go-kit logger for logging req and logic error
func createLogger(serviceName string) log.Logger {
	// initialize logger and set it to logger in global scope
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(
		logger,
		"service", serviceName,
		"time", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	// logging that the service has started
	level.Info(logger).Log(MessageLog, "service started")
	// logging that the service has ended
	defer level.Info(logger).Log(MessageLog, "service ended")
	return logger
}

func openDatabase() *gorm.DB {

	db, err := dbconnector.DBCredential{
		DBDriver:     os.Getenv("DB_DRIVER"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBSqlitePath: "",
	}.Connect()
	if err != nil {
		fmt.Println(err)
	}

	// Load debuging mode env
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	db.LogMode(debug)

	return db
}

func createMicroAuthService(db *gorm.DB) svc.Service {
	microAuthRepository := mysql.NewDBReadWriter(db)
	return usecase.NewMicroAuthService(microAuthRepository)
}

func createEndpoint(service svc.Service) transport.Endpoints {
	// create endpoints instance
	endpoints := transport.MakeEndpoints(service)

	// attach logging middleware into every endpoints
	endpoints.Login = middleware.LoggingMiddleware(log.With(logger, "method", "Login"))(endpoints.Login)

	// attach circuit breaker middleware into every endpoints
	endpoints.Login = middleware.CircuitBreakerMiddleware(microAuthHystrixCommand)(endpoints.Login)

	return endpoints
}

func grpcGatewayMode(
	ctx context.Context,
	hierarchyServiceGrpc pb.AuthServiceServer,
) {
	port := os.Getenv("GATEWAY_PORT")
	mux := runtime.NewServeMux()
	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, hierarchyServiceGrpc)
	if err != nil {
		_ = level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	listen, err := net.Listen("tcp", port)
	if err != nil {
		_ = level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		_ = logger.Log("transport", "gRPC and Restful", "addr", port)
		errs <- http.Serve(listen, mux)
	}()
	_ = level.Error(logger).Log("exit", <-errs)
}

func configureHystrixCommand(command string) {
	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})
}

func App() {
	// configure hystrix
	configureHystrixCommand(microAuthHystrixCommand)
	// creAuthServiceServerate context
	ctx := context.Background()
	// create logger
	logger = createLogger(serviceName)
	// load .env file content
	loadEnv()
	// open database
	db := openDatabase()
	// prepare service
	service := createMicroAuthService(db)
	// prepare endpoints
	endpoints := createEndpoint(service)
	// create hierarchy service grpc
	hierarchyServiceGrpc := transportgrpc.NewGRPCServer(endpoints, logger)
	// start grpc server
	//grpcMode(ctx, hierarchyServiceGrpc)
	grpcGatewayMode(ctx, hierarchyServiceGrpc)
	// clear context
	defer ctx.Done()
}
