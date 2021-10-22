package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/user_srv/proto"
)

var (
	userClient proto.UserClient
	conn       *grpc.ClientConn
)

func TestGetUserList() {
	ctx := context.Background()
	list, err := userClient.GetUserList(ctx, &proto.PageInfo{
		Pn:    0,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range list.Data {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		check, err := userClient.CheckPassword(ctx, &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(check.Success)
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func main() {
	Init()
	TestGetUserList()
	conn.Close()
}
