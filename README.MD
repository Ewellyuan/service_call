### 本实例基于 go1.10.3

* 监听端口：9090

* 安装依赖
```
RUN go get -d -v ./...
RUN go install -v ./...
```

* 运行
```
go run main.go
```

* 打包镜像
```
docker build -t service_call_v1 .
```