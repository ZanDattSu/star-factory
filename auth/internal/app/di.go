package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ZanDattSu/star-factory/auth/internal/api/v1/auth"
	"github.com/ZanDattSu/star-factory/auth/internal/api/v1/user"
	"github.com/ZanDattSu/star-factory/auth/internal/config"
	"github.com/ZanDattSu/star-factory/auth/internal/repository"
	sessionRepo "github.com/ZanDattSu/star-factory/auth/internal/repository/session"
	userRepo "github.com/ZanDattSu/star-factory/auth/internal/repository/user"
	"github.com/ZanDattSu/star-factory/auth/internal/service"
	authService "github.com/ZanDattSu/star-factory/auth/internal/service/auth"
	"github.com/ZanDattSu/star-factory/auth/internal/service/hasher"
	userService "github.com/ZanDattSu/star-factory/auth/internal/service/user"
	"github.com/ZanDattSu/star-factory/platform/pkg/cache"
	rediscache "github.com/ZanDattSu/star-factory/platform/pkg/cache/redis"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authApi authV1.AuthServiceServer
	userApi userV1.UserServiceServer

	authService    service.AuthService
	userService    service.UserService
	passwordHasher service.PasswordHasher

	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository

	redisClient    cache.RedisClient
	redisPool      *redigo.Pool
	postgreSQLPool *pgxpool.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthApi(ctx context.Context) authV1.AuthServiceServer {
	if d.authApi == nil {
		d.authApi = auth.NewAuthApi(
			d.AuthService(ctx),
		)
	}

	return d.authApi
}

func (d *diContainer) UserApi(ctx context.Context) userV1.UserServiceServer {
	if d.userApi == nil {
		d.userApi = user.NewUserApi(
			d.UserService(ctx),
		)
	}

	return d.userApi
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewAuthService(
			d.UsersRepository(ctx),
			d.SessionsRepository(),
			config.AppConfig().Session.TTL(),
			d.PasswordHasher(),
		)
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewUsersService(
			d.UsersRepository(ctx),
			d.PasswordHasher(),
		)
	}

	return d.userService
}

func (d *diContainer) PasswordHasher() service.PasswordHasher {
	if d.passwordHasher == nil {
		d.passwordHasher = hasher.NewBcryptHasher()
	}

	return d.passwordHasher
}

func (d *diContainer) UsersRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepo.NewUserRepository(
			d.PostgreSQLPool(ctx),
		)
	}

	return d.userRepository
}

func (d *diContainer) SessionsRepository() repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepo.NewSessionRepository(
			d.RedisClient(),
		)
	}

	return d.sessionRepository
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = rediscache.NewClient(
			d.RedisPool(),
			logger.Logger(),
			config.AppConfig().Redis.ConnectionTimeout(),
		)
	}

	return d.redisClient
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgreSQLPool == nil {
		dbURI := config.AppConfig().Postgres.URI()

		pool, err := pgxpool.New(ctx, dbURI)
		if err != nil {
			panic(fmt.Sprintf("Failed to create pgxpool connect: %s", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("Database is unavailable: %s", err))
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgreSQLPool = pool
	}

	return d.postgreSQLPool
}
