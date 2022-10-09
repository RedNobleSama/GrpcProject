# GrpcProject
GRPC学习以及一些中间件使用

# ProtoBUf
goland ide 添加protocbuf文件路径

protoc -I . base.proto user.proto --go_out=plugins=grpc:.
protoc -I .  user.proto --go_out=plugins=grpc:.

引用protobuf自带proto
D:\develop\program\Go\GOPATH\src 为proto文件所在路径
protoc -I . -I D:\develop\program\Go\GOPATH\src  user.proto --go_out=plugins=grpc:.

linux下使用proto-gen-go 操作proto文件生成go代码
go get github.com/golang/protobuf/protoc-gen-go
cd github.com/golang/protobuf/protoc-gen-go
go build
将bin/protoc-gen-go加入环境变量

将下载的对应goprotobuf支持库 放到gopath/src下面

CN:16