package response

import (
	"api/errorx"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Success(w, r, map[string]string{"foo": "bar"})

	var resp Body
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	if resp.Code != errorx.OK {
		t.Fatalf("期望 code %d, 得到 %d", errorx.OK, resp.Code)
	}
	if resp.Msg != "success" {
		t.Fatalf("期望 msg 'success', 得到 '%s'", resp.Msg)
	}
	if resp.Data == nil {
		t.Fatal("data 不应为 nil")
	}
}

func TestError_WithCodeErr(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	codeErr := errorx.NewCodeError(errorx.ErrUserNotFound)
	Error(w, r, codeErr)

	var resp Body
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	if resp.Code != errorx.ErrUserNotFound {
		t.Fatalf("期望 code %d, 得到 %d", errorx.ErrUserNotFound, resp.Code)
	}
	if resp.Msg != "用户不存在" {
		t.Fatalf("期望 msg '用户不存在', 得到 '%s'", resp.Msg)
	}
}

func TestError_WithGenericError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Error(w, r, nil)

	var resp Body
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	if resp.Code != errorx.ErrUnknown {
		t.Fatalf("期望 code %d, 得到 %d", errorx.ErrUnknown, resp.Code)
	}
	if resp.Msg != "内部错误" {
		t.Fatalf("期望 msg '内部错误', 得到 '%s'", resp.Msg)
	}
}
