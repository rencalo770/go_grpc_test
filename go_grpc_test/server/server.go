package main

import (
	"fmt"
	proto "go_grpc_test/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)


var inputFile string

func init() {

	//todo 在这里输入要发送的文件路径
	inputFile = ""

	if inputFile == "" {
		panic("请输入要发送的文件名!")
	}
}


type Server struct {
	proto.UnimplementedFileServiceServer
}

func (s *Server) GetFile(request *proto.FileRequest, stream proto.FileService_GetFileServer) error{

	println("rec-->", request.GetId())
	fileRead, err := os.OpenFile(inputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	for  {
		var bs = make([]byte, 1000)
		i, err := fileRead.Read(bs)
		if err != nil {
			println( "文件读取完毕:", fmt.Sprintf("%v", err))
			break
		}

		if i > 0 {

			bsx := make([]byte, i)
			for j := 0; j < i; j++ {
				bsx[j] = bs[j]
			}

			err := stream.Send(&proto.FileResponse{
				Length: int32(i),
				Bs:     bsx,
			})
			if err != nil {
				println(fmt.Sprintf("发送文件失败:%v", err))
				break
			}

			//重置
			i = 0
		}else {
			println("文件读取完毕")
			break
		}
	}

	_ = fileRead.Close()

	//这条流写完毕
	stream.Context().Done()
	return nil
}

//todo 注意不要去实现这个方法
//func (s *Server) mustEmbedUnimplementedFileServiceServer(){
//	return
//}

//var _ pb.FileServiceServer = &Server{}





func main() {

	// 监听127.0.0.1:50051地址
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc服务端
	s := grpc.NewServer()

	server2 := &Server{}
	// 注册服务
	proto.RegisterFileServiceServer(s, server2)

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
