# GrpcProject
GRPC学习以及一些中间件使用

# ProtoBUf
goland ide 添加protocbuf文件路径

protoc -I . base.proto user.proto --go_out=plugins=grpc:.
protoc -I .  user.proto --go_out=plugins=grpc:.

引用protobuf自带proto
D:\develop\program\Go\GOPATH\src 为proto文件所在路径
protoc -I . -I D:\develop\program\Go\GOPATH\src  user.proto --go_out=plugins=grpc:.