package logic

import (
	"context"
	"time"

	"user/internal/svc"
	"user/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	token, err := jwt.Parse(in.RefreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, "无效的 token")
		}
		return []byte(l.svcCtx.Config.JwtSecret), nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "refresh token 无效或已过期")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "无效的 token 声明")
	}

	userID := int64(claims["user_id"].(float64))
	now := time.Now().Unix()

	newAccessToken, err := l.generateToken(userID, now, 3600)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成 token 失败")
	}
	newRefreshToken, err := l.generateToken(userID, now, 604800)
	if err != nil {
		return nil, status.Error(codes.Internal, "生成 refresh token 失败")
	}

	return &user.RefreshTokenResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (l *RefreshTokenLogic) generateToken(userId int64, now int64, ttl int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["iat"] = now
	claims["exp"] = now + ttl
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.JwtSecret))
}
