package logic

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.UpdateUserResponse, err error) {
	_, err = l.svcCtx.UserRpc.UpdateUser(l.ctx, &user.UpdateUserRequest{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateUserResponse{}, nil
}
