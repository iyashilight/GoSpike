package logic

import (
	"context"
	"time"

	"user/internal/svc"
	"user/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx,in.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound,"用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(in.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated,"密码错误")
	}
	now := time.Now().Unix()
	accessToken, err := l.generateToken(u.Id, now, 3600)
	if err != nil {
		return nil, err
	}
	refreshToken, err := l.generateToken(u.Id, now, 604800) // 7 days
	if err != nil {
		return nil, err
	}
	return &user.LoginResponse{
		Id:           u.Id,
		Name:         u.Name,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (l *LoginLogic) generateToken(userId int64, now int64, ttl int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["iat"] = now
	claims["exp"] = now + ttl
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.JwtSecret))
}
