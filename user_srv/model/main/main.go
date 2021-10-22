package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"strings"
)

func main() {
	// Using custom options
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(salt, newPassword)
	passwordInfo := strings.Split(newPassword, "$")
	fmt.Println(passwordInfo)
	check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // true

	// 插入测试数据
	for i := 0; i < 10; i++ {
		user := model.User{
			Mobile:   fmt.Sprintf("1391301234%d", i),
			Password: newPassword,
			NickName: fmt.Sprintf("bobby%d", i),
		}
		global.DB.Save(&user)
	}
}
