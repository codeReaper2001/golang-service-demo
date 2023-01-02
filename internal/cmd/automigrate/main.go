package main

import (
	"context"
	"fmt"
	"go_test/internal/config"
	"go_test/pkg/ent"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlConfig = config.MysqlConfig{
	Username: "root",
	Password: "123456",
	Host:     "192.168.0.100",
	Port:     "3340",
	DBName:   "dev",
}

func main() {
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
