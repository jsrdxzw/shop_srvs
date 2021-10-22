package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
	"strings"
	"time"
)

type UserService struct {
	*proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u *UserService) GetUserList(ctx context.Context, info *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	var total int64
	result := global.DB.Model(&model.User{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(total)
	// 分页
	global.DB.Scopes(Paginate(int(info.Pn), int(info.PSize))).Find(&users)
	for _, user := range users {
		userInfoResponse := Model2Response(user)
		rsp.Data = append(rsp.Data, userInfoResponse)
	}
	return rsp, nil
}

func (u *UserService) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: request.Mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoResp := Model2Response(user)
	return userInfoResp, nil
}

func (u *UserService) GetUserById(ctx context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	// 通过ID查询用户
	result := global.DB.First(&user, request.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoResp := Model2Response(user)
	return userInfoResp, nil
}

func (u *UserService) CreateUser(ctx context.Context, info *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: info.Mobile}).First(&user)
	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = info.Mobile
	user.NickName = info.NickName
	// 密码加密
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(info.Password, options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.Password = newPassword
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	response := Model2Response(user)
	return response, nil
}

func (u *UserService) UpdateUser(ctx context.Context, info *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, info.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(info.Birthday), 0)
	user.NickName = info.NickName
	user.Birthday = &birthday
	user.Gender = info.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserService) CheckPassword(ctx context.Context, info *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(info.EncryptedPassword, "$")
	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}

func Model2Response(user model.User) *proto.UserInfoResponse {
	userInfoResp := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     uint32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix())
	}
	return &userInfoResp
}
