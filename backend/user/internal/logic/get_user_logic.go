package logic

import (
	"context"
	"time"

	"user/internal/svc"
	"user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx,in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound,"用户不存在")
	}
	return &user.GetUserResponse{
		Id: u.Id,
		Name: u.Name,
		Email: u.Email,
		CreatedAt: u.CreatedAt.Format(time.DateTime),
	}, nil
}
