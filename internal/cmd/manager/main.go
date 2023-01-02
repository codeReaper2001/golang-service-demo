package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	student_pb "go_test/api/student"
	teacher_pb "go_test/api/teacher"
	"go_test/internal/config"
	student_svc "go_test/internal/service/student"
	teacher_svc "go_test/internal/service/teacher"
	"go_test/pkg/ent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
}

type endPointFunction func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

var endPointFunctions = []endPointFunction{
	student_pb.RegisterStudentSvcHandlerFromEndpoint,
	teacher_pb.RegisterTeacherSvcHandlerFromEndpoint,
}

func run(conf *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	entClient := NewEntClient(&conf.MySQL)
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcServer := grpc.NewServer()
	grpcServerAddr := fmt.Sprintf("localhost:%s", conf.Server.GrpcServer)
	lis, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		logger.Fatalf("bind err: %v", err)
	}
	student_pb.RegisterStudentSvcServer(grpcServer, student_svc.New(logger, entClient))
	teacher_pb.RegisterTeacherSvcServer(grpcServer, teacher_svc.New(logger, entClient))
	go grpcServer.Serve(lis)
	logger.Infoln("grpcServer started")
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	for _, endPointFunc := range endPointFunctions {
		endPointFunc(ctx, mux, grpcServerAddr, opts)
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(fmt.Sprintf(":%s", conf.Server.GrpcGateway), mux)
}

func loadConfigFile(filePath string) *config.Config {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config err: %v\n", err)
	}
	logger.Infoln("read config ok")
	conf := &config.Config{}
	viper.Unmarshal(conf)
	logger.Infoln(conf)
	return conf
}

func NewEntClient(dbConfig *config.MysqlConfig) *ent.Client {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	logger.Infoln("connect to db ok!")
	return client
}

func main() {
	flag.Parse()
	defer glog.Flush()

	// 加载配置文件
	conf := loadConfigFile(os.Args[1])

	if err := run(conf); err != nil {
		glog.Fatal(err)
	}
}
