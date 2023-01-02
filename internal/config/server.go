package config

type ServerConfig struct {
	GrpcGateway string `mapstructure:"grpc_gateway"`
	GrpcServer  string `mapstructure:"grpc_server"`
}
