package logic

import (
	"context"
	"errors"
	"unicode"

	"user/internal/model"
	"user/internal/svc"
	"user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	existing, _ := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if existing != nil {
		return nil, status.Error(codes.AlreadyExists, "邮箱已被注册")
	}

	if err := validatePassword(in.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "密码加密失败")
	}
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: string(hashedPwd),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "注册失败")
	}
	id, _ := result.LastInsertId()
	return &user.RegisterResponse{Id: id}, nil
}

func validatePassword(pwd string) error {
	if len(pwd) < 8 {
		return errors.New("密码长度至少 8 位")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
	)

	for _, ch := range pwd {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.New("密码必须包含大写字母")
	}
	if !hasLower {
		return errors.New("密码必须包含小写字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含数字")
	}
	return nil
}
