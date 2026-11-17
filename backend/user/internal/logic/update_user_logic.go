package logic

import (
	"context"

	"user/internal/model"
	"user/internal/svc"
	"user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		return nil, status.Error(codes.Internal, "查询用户失败")
	}

	name := in.Name
	if name == "" {
		name = u.Name
	}
	email := in.Email
	if email == "" {
		email = u.Email
	}

	// 保持原有密码不变，Update 方法包含 password 字段
	err = l.svcCtx.UserModel.Update(l.ctx, &model.User{
		Id:       in.Id,
		Name:     name,
		Email:    email,
		Password: u.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "更新用户失败")
	}

	return &user.UpdateUserResponse{}, nil
}
