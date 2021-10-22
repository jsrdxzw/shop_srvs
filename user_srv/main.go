package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"shop_srvs/user_srv/handler"
	"shop_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})
	url := fmt.Sprintf("%s:%d", *IP, *Port)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatalf("failed to listen: " + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start grpc: " + err.Error())
	}

}
