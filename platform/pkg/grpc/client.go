package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
}
