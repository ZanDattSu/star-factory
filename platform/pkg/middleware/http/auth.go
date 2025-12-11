package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// AuthClient это алиас для сгенерированного gRPC клиента
type AuthClient = authV1.AuthServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	authClient AuthClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(authClient AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session UUID из заголовка
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.authClient.Whoami(r.Context(), &authV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		ctx = grpcAuth.AddUserToContext(ctx, whoamiRes.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
