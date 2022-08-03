package main

import (
	"context"
	"fmt"
	proto "go_grpc_test/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
)

var outputFile string

func init() {
	//todo 在这里输入要存储的文件路径
	outputFile = ""

	if outputFile == "" {
		panic("请输入要存储的文件名!")
	}
}

func main() {

	// 1.建立连接 获取client
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewFileServiceClient(conn)
	// 2.执行各个Stream的对应方法
	stream, err := client.GetFile(context.Background(), &proto.FileRequest{Id: 1})
	if err != nil {
		log.Printf("send error:%v\n", err)
	}


	fileWrite, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	for  {
		recv, err := stream.Recv()
		if err == io.EOF {
			log.Printf("server closed: %v", err)
			break
		}

		if err != nil {
			log.Printf("Recv error:%v", err)
			continue
		}

		println("len----->", recv.Length)

		i, err := fileWrite.Write(recv.Bs)
		if err != nil {
			println(fmt.Sprintf("写文件失败:%v", err))
			break
		}
		println("写入文件长度:", i)
	}

	_ = fileWrite.Close()
}
