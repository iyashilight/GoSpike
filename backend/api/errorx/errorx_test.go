package errorx

import (
	"testing"
)

func TestNewCodeError_KnownCode(t *testing.T) {
	err := NewCodeError(ErrUserNotFound)
	if err.Code != ErrUserNotFound {
		t.Fatalf("期望 code %d, 得到 %d", ErrUserNotFound, err.Code)
	}
	if err.Msg != "用户不存在" {
		t.Fatalf("期望 msg '用户不存在', 得到 '%s'", err.Msg)
	}
}

func TestNewCodeError_UnknownCode(t *testing.T) {
	err := NewCodeError(99999)
	if err.Code != ErrUnknown {
		t.Fatalf("未知 code 应 fallback 到 ErrUnknown(%d), 得到 %d", ErrUnknown, err.Code)
	}
}

func TestNewCodeError_CustomMsg(t *testing.T) {
	err := NewCodeError(ErrStockNotEnough, "自定义消息")
	if err.Code != ErrStockNotEnough {
		t.Fatalf("期望 code %d, 得到 %d", ErrStockNotEnough, err.Code)
	}
	if err.Msg != "自定义消息" {
		t.Fatalf("期望 msg '自定义消息', 得到 '%s'", err.Msg)
	}
}

func TestNewCodeError_AllCodesHaveMessages(t *testing.T) {
	codes := []int{
		OK, ErrUserNotFound, ErrInvalidCredential, ErrUserAlreadyExists,
		ErrProductNotFound, ErrStockNotEnough, ErrCreateProductFailed,
		ErrUpdateProductFailed, ErrDeleteProductFailed, ErrCacheStockFailed,
		ErrOrderNotFound, ErrCreateOrderFailed,
	}
	for _, code := range codes {
		err := NewCodeError(code)
		if err.Msg == "" {
			t.Errorf("code %d 缺少错误消息", code)
		}
	}
}

func TestCodeErr_Error(t *testing.T) {
	err := NewCodeError(ErrStockNotEnough)
	if err.Error() != "库存不足" {
		t.Fatalf("Error() 应返回消息内容, 得到 '%s'", err.Error())
	}
}
