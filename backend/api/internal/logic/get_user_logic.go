// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"user/user"

	"api/errorx"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserRequest{
		Id: req.Id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.NotFound {
			return nil, errorx.NewCodeError(errorx.ErrUserNotFound)
		}
		return nil, err
	}
	return &types.GetUserResponse{
		Id:    rpcResp.Id,
		Name:  rpcResp.Name,
		Email: rpcResp.Email,
	}, nil
}
