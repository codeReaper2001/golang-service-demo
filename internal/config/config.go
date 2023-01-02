package config

type Config struct {
	MySQL  MysqlConfig  `mapstructure:"mysql"`
	Server ServerConfig `mapstructure:"server"`
}
