# protoc-gen-apidoc
根据 Proto 文件生成接口文档。支持 Swagger、Postman、HTML、Markdown 等多种格式输出。

---



### 安装

- ##### [protoc](https://github.com/protocolbuffers/protobuf/releases)

- ##### protoc--gen-go

  ```shell
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

  # protoc-gen-go-grpc
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

- ##### protoc--gen-gogo

  ```shell
  go install github.com/gogo/protobuf/protoc-gen-gogo@latest
  ```

- ##### protoc-gen-apidoc

  ```shell
  # 方式一
  go install github.com/charlesbases/protoc-gen-apidoc@latest

  # 方式二
  git clone https://github.com/charlesbases/protoc-gen-apidoc.git
  cd protoc-gen-apidoc && go install .
  ```

- ##### google/protobuf/*.proto

  ```shell
  git clone https://github.com/charlesbases/protobuf.git
  cp -r protobuf/google ${GOPATH}/src/.
  # 或
  cd protobuf && make init
  ```

### 运行

```shell
# 默认配置文件路径(./aidoc.yaml)
protoc -I=${GOPATH}/src:. --apidoc_out=. pb/*.proto

# 指定配置文件路径
protoc -I=${GOPATH}/src:. --apidoc_out=configfile=others/apidoc.yaml:. pb/*.proto
```

### 参数

- ##### configfile: 配置文件路径。默认 "./apidoc.yaml"

### 配置文件

```yaml
# api host
host: 127.0.0.1
# service name
title: test
# 添加 header
header:
  - Authorization
# 生成文档类型
document:
  # HTML
  - type: html
    # 文档路径
    file: api.html
  # Markdown
  - type: markdown
    file: api.md
  # Postman
  - type: postman
    file: postman.json
  # Swagger API
  - type: swagger
    file: swagger.json

```

### proto 文件注释格式

- ##### 格式一: 默认请求方式为 POST

  ```protobuf
  syntax = "proto3";

  option go_package = ".;pb";

  package pb;

  import "pb/base.proto";

  // 用户服务
  service Users {
    // 用户列表
    rpc List (Request) returns (Response) {}
  }

  // 入参
  message Request {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }

  // 出参
  message Response {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  ```

- ##### 格式二: 自定义请求方式、请求路径、Content-Type

  ```protobuf
  syntax = "proto3";
  
  option go_package = ".;pb";
  
  package pb;
  
  import "pb/base.proto";
  import "google/protobuf/plugin/http.proto";
  
  // 用户服务
  service Users {
    // 获取用户
    rpc User (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users/{uid}"
      };
    }
    // 用户列表
    rpc UserList (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users"
      };
    }
    // 创建用户
    rpc UserCreate (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        post: "/api/v1/users"
      };
    }
    // 更新用户
    rpc UserUpdate (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        put: "/api/v1/users/{uid}"
      };
    }
    // 删除用户
    rpc UserDelete (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        delete: "/api/v1/users/{uid}"
      };
    }
    // 用户头像上传
    rpc UserUpload (Upload) returns (Response) {
      option (google.protobuf.plugin.http) = {
        put: "/api/v1/users/{uid}"
        consume: "multipart/form-data"
      };
    }
    // 用户头像下载
    rpc UserUpload (Request) returns (Upload) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users/{uid}"
        produce: "multipart/form-data"
      };
    }
  }
  
  // 入参
  message Request {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 出参
  message Response {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 头像上传
  message Upload {
   FileType type = 1;
   bytes file = 2;
  }
  
  // 图片类型
  enum FileType {
    JPG = 0;
    PNG = 1;
    GIF = 2;
  }
  ```
