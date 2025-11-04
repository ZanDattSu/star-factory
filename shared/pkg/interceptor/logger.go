package interceptor

import (
	"context"
	"log"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor создает серверный унарный интерцептор, который логирует
// информацию о времени выполнения методов gRPC сервера.
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем имя метода из полного пути
		method := path.Base(info.FullMethod)

		log.Printf("🚀 Started gRPC method %s\n", method)

		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
		}

		return resp, err
	}
}
