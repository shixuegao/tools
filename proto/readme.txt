关于google protobuffer的使用
1.下载编译器protoc, 网址https://github.com/google/protobuf/releases  
2.下载protoc在go下运行的插件, go get github.com/golang/protobuf
3.使用protoc --go-grpc_out=. .\proto\route_guide.proto来生成go文件
