module payment

go 1.24.0

replace github.com/ZanDattSu/star-factory/shared => ../shared

require (
	github.com/ZanDattSu/star-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250929231259-57b25ae835d4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251103181224-f26f9409b101 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
