/**
    @auther: oreki
    @date: 2022/4/25
    @note: 图灵老祖保佑,永无BUG
**/

package controller

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
	"user_srv/db"
	in "user_srv/interface"
	"user_srv/model"
)

type UserServer struct {
}

// HashPassword 密码加密
func HashPassword(password []byte) []byte {
	hashPWS, _ := bcrypt.GenerateFromPassword(password, 10)
	return hashPWS
}

// ComparePassword 密码校验
func ComparePassword(hashedPassword string, password string) bool {
	byteHash := []byte(hashedPassword)
	bytePassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePassword)
	if err != nil {
		fmt.Println("密码错误", err.Error())
		return false
	}
	return true
}

// ModelToResponse 格式转换
func ModelToResponse(user model.User) *in.UserInfoResponse {
	// 在grpc的message中字段有默认值，不能随便赋值nil进去
	// 这里要搞清楚哪些字段有默认值
	userInfoRes := &in.UserInfoResponse{
		Id:       int64(user.ID),
		NickName: user.NickName,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     int64(user.Role),
	}
	if user.Birthday != nil {
		userInfoRes.Birthday = user.Birthday.Unix()
	}
	return userInfoRes
}

// Paginate gorm 通用分页
func Paginate(pageSize, pageNumber int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNumber == 0 {
			pageNumber = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNumber - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context, req *in.PageInfo) (*in.UserListResponse, error) {
	var users []model.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &in.UserListResponse{
		Total: result.RowsAffected,
	}

	db.DB.Scopes(Paginate(int(req.PSize), int(req.PNum))).Find(&users)
	for _, user := range users {
		rsp.Data = append(rsp.Data, ModelToResponse(user))
	}

	return rsp, nil
}

// GetUserbyMobile 获取用户信息 by mobile
func (s *UserServer) GetUserbyMobile(ctx context.Context, req *in.MobileRequest) (*in.UserInfoResponse, error) {
	var user model.User
	result := db.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

// GetUserbyId 获取用户信息 by id
func (s *UserServer) GetUserbyId(ctx context.Context, req *in.IdRequest) (*in.UserInfoResponse, error) {
	var user model.User
	result := db.DB.Where("id = ?", req.Id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *in.CreateUserInfo) (*in.UserInfoResponse, error) {
	var user model.User
	result := db.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已注册")
	}

	user = model.User{
		Mobile:   req.Mobile,
		NickName: req.NickName,
		Password: string(HashPassword([]byte(req.Password))),
	}
	result = db.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error()) // 内部错误
	}

	rsp := ModelToResponse(user)
	return rsp, nil
}

// UpdateUser 修改用户
func (s *UserServer) UpdateUser(ctx context.Context, req *in.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := db.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	Brithday := time.Unix(req.Birthday, 0)
	result = db.DB.Updates(&model.User{
		NickName: req.NickName,
		Birthday: &Brithday,
		Gender:   req.Gender,
	})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error()) // 内部错误
	}
	return &empty.Empty{}, nil
}

// CheckUserPassword 检查用户密码
func (s *UserServer) CheckUserPassword(ctx context.Context, req *in.PasswordInfo) (*in.CheckResponse, error) {
	result := ComparePassword(req.EncrpytedPassword, req.Password)
	rsp := in.CheckResponse{
		Success: result,
	}
	return &rsp, nil
}
