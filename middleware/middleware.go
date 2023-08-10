package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/egiferdians/micro-auth/util/errcode"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/codes"
)

// CircuitBreakerMiddleware for endpoint
func CircuitBreakerMiddleware(command string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			var logicErr error
			err = hystrix.Do(command, func() (err error) {
				resp, logicErr = next(ctx, request)
				return logicErr
			}, func(err error) error {
				return err
			})
			if logicErr != nil {
				return nil, logicErr
			}
			if err != nil {
				return nil, errcode.New(
					codes.Unavailable,
					errors.New("service is busy or unavailable, please try again later"),
				)
			}
			return resp, nil
		}
	}
}

// LoggingMiddleware for endpoint
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var resp interface{}
			req, _ := json.Marshal(request)
			defer func(begin time.Time) {
				logger.Log(
					"transport_error", err,
					"took", time.Since(begin),
					"request", string(req),
				)
			}(time.Now())
			resp, err = next(ctx, request)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}
