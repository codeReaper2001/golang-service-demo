## ent相关

### 生成原始ent文件

执行：`go run -mod=mod entgo.io/ent/cmd/ent init User`，生成相关文件

修改schema目录中的内容后，执行

```bash
go generate ./path_to_ent/ent
```

即可生成相关文件

### 自动迁移到数据库代码

```go
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
```

运行此main.go即可将变更迁移至数据库

## protobuf相关

### 1、下载buf

bash脚本：

```bash
# Substitute BIN for your bin directory.
# Substitute VERSION for the current released version.
BIN="/usr/local/bin" && \
VERSION="1.11.0" && \
  curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
    -o "${BIN}/buf" && \
  chmod +x "${BIN}/buf"
```

### 2、配置

文档：https://docs.buf.build/tour/introduction

下载go相关代码生成插件：https://docs.buf.build/tour/generate-go-code

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

编写 buf.gen.yaml：

```yaml
version: v1
plugins:
  - name: go
    out: ./
    opt:
      - paths=source_relative
  - name: go-grpc
    out: ./
    opt:
      - paths=source_relative
+  - name: grpc-gateway
+    out: ./
+    opt:
+      - paths=source_relative
+      - generate_unbound_methods=true
```

放置proto目录及文件，如：

```bash
api/
├── buf.gen.yaml
├── buf.work.yaml
└── student_apis
    ├── buf.yaml
    └── student
        ├── rpc.proto
        └── student.proto
```

其中buf.yaml内容如下，它是使用buf init命令生成的：

```yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
+ deps:
+  - buf.build/googleapis/googleapis
```

student.proto：

```protobuf
syntax = "proto3";

package student;

option go_package = "go_test/student";

message Student {
    int32 stu_id = 1;
    string name = 2;
    int32 age = 3;
}

message GetStudentRequest {
    int32 stu_id = 1;
}

message GetStudentResponse {
    int32 status = 1;
    string msg = 2;
    Student data = 3;
}
```

rpc.proto：

```protobuf
syntax = "proto3";

package student;
option go_package = "go_test/student";

import "student/student.proto";
import "google/api/annotations.proto";

service StudentSvc {
    rpc GetStudent(GetStudentRequest) returns (GetStudentResponse) {
        option (google.api.http) = {
            get: "/student"
        };
    }
}
```

这里会依赖 `annotations.proto` ，需要进入`student_apis/`，执行`buf mod update`，生成一个 `buf.lock` 文件，即可使用

然后在外层添加 `buf.gen.yaml` 和 `buf.work.yaml` ，内容如下：

（grpc-gateway相关protoc插件下载参考：https://github.com/grpc-ecosystem/grpc-gateway）

```yaml
# buf.gen.yaml
version: v1
plugins:
  - name: go
    out: .
    opt:
      - paths=source_relative
  - name: go-grpc
    out: .
    opt:
      - paths=source_relative
  - name: grpc-gateway
    out: .
    opt:
      - paths=source_relative
      - generate_unbound_methods=true
  - name: openapiv2
    out: .
```

```yaml
# buf.work.yaml
version: v1
directories:
  - student_apis
```

最后在`api`目录下执行：

```bash
buf generate
```

即可生成 `student` 和 `grpc-gateway` 相关的pb文件和swagger文件：

```bash
api/student
├── rpc_grpc.pb.go
├── rpc.pb.go
├── rpc.pb.gw.go
├── rpc.swagger.json
├── student.pb.go
└── student.swagger.json
```

## 实现student服务

```go
package student

import (
    "context"

    student_pb "go_test/api/student"
)

type service struct {
}

func New() *service {
    return &service{}
}

func (s *service) GetStudent(
    ctx context.Context,
    req *student_pb.GetStudentRequest,
) (*student_pb.GetStudentResponse, error) {
    return &student_pb.GetStudentResponse{
        Status: 0,
        Msg:    "OK",
        Data:   &student_pb.Student{},
    }, nil
}
```

## 编写入口main.go文件

参考链接：https://github.com/grpc-ecosystem/grpc-gateway

```go
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
    student_svc "go_test/internal/service/student"

    _ "github.com/go-sql-driver/mysql"
)

type endPointFunction func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

var endPointFunctions = []endPointFunction{
    student_pb.RegisterStudentSvcHandlerFromEndpoint,
}

func run(conf *config.Config) error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    // grpc服务
    grpcServer := grpc.NewServer()
    grpcServerAddr := "localhost:10001"
    lis, err := net.Listen("tcp", grpcServerAddr)
    if err != nil {
        logger.Fatalf("bind err: %v", err)
    }
    student_pb.RegisterStudentSvcServer(grpcServer, student_svc.New(logger, entClient))
    teacher_pb.RegisterTeacherSvcServer(grpcServer, teacher_svc.New(logger, entClient))
    // grpc启动
    go grpcServer.Serve(lis)
    logger.Infoln("grpcServer started")
    // Register gRPC server endpoint
    // Note: Make sure the gRPC server is running properly and accessible
    // 配置grpc-gateway对各服务的转发规则
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    for _, endPointFunc := range endPointFunctions {
        endPointFunc(ctx, mux, grpcServerAddr, opts)
    }
    // Start HTTP server (and proxy calls to gRPC server endpoint)
    // 启动grpc-gateway服务
    return http.ListenAndServe(fmt.Sprintf(":10000", mux)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if err := run(conf); err != nil {
        glog.Fatal(err)
    }
}

```

这样即可运行后端程序

## 读取配置文件

第三方库viper：`"github.com/spf13/viper"`

```go
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
```

相关配置文件：

```yaml
mysql:
  username: "root"
  password: "123456"
  host: "192.168.0.100"
  port: "3340"
  db_name: "dev"
```

映射到的结构体：

```go
type Config struct {
    MySQL MysqlConfig `mapstructure:"mysql"`
}
type MysqlConfig struct {
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    DBName   string `mapstructure:"db_name"`
}
```

## 日志打印配置

第三方库：`"github.com/sirupsen/logrus"`

配置：

```go
var logger *logrus.Logger

func init() {
    logger = logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{
        ForceQuote:      true,                  //键值对加引号
        TimestampFormat: "2006-01-02 15:04:05", //时间格式
        FullTimestamp:   true,
    })
}
```

