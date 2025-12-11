package app

import (
	"context"
	"fmt"

	payApi "github.com/ZanDattSu/star-factory/payment/internal/api/v1/payment"
	"github.com/ZanDattSu/star-factory/payment/internal/config"
	"github.com/ZanDattSu/star-factory/payment/internal/service"
	payService "github.com/ZanDattSu/star-factory/payment/internal/service/payment"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	grpcclient "github.com/ZanDattSu/star-factory/platform/pkg/grpc"
	"github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1Api   paymentV1.PaymentServiceServer
	paymentService service.PaymentService

	authClient      authV1.AuthServiceClient
	authInterceptor *interceptor.AuthInterceptor
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1Api(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1Api == nil {
		d.paymentV1Api = payApi.NewApi(d.PaymentService(ctx))
	}

	return d.paymentV1Api
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payService.NewService()
	}

	return d.paymentService
}

func (d *diContainer) AuthClient(_ context.Context) authV1.AuthServiceClient {
	if d.authClient == nil {
		authConn, err := grpcclient.NewGRPCConnectWithoutSecure(config.AppConfig().Auth.AuthServiceAddress())
		if err != nil {
			panic(fmt.Sprintf(
				"Failed to connect to Auth gRPC service (%s): %v",
				config.AppConfig().Auth.AuthServicePort(),
				err,
			))
		}

		closer.AddNamed("Auth connection", func(ctx context.Context) error {
			return authConn.Close()
		})

		d.authClient = authV1.NewAuthServiceClient(authConn)
	}

	return d.authClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *interceptor.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = interceptor.NewAuthInterceptor(d.AuthClient(ctx))
	}

	return d.authInterceptor
}
