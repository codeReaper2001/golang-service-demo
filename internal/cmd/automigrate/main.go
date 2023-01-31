package main

import (
	"context"
	"fmt"
	"go_test/internal/config"
	"go_test/pkg/ent"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfigFile(filePath string) *config.Config {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config err: %v\n", err)
	}
	logrus.Infoln("read config ok")
	conf := &config.Config{}
	viper.Unmarshal(conf)
	logrus.Infoln(conf)
	return conf
}

func main() {
	config := loadConfigFile("./config.yaml")
	mysqlConfig := config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.DBName,
	)
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Successfully Migrated!")
}
