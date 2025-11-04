package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validatable interface {
	Validate() error
}

// ValidationInterceptor LoggerInterceptor создает серверный унарный интерцептор, который
// валидирует входящие запросы
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if v, ok := req.(Validatable); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("validation failed: %v", err))
			}
		}

		return handler(ctx, req)
	}
}
