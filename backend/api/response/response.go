package response

import (
	"api/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(w http.ResponseWriter, r *http.Request, data any) {
	httpx.OkJsonCtx(r.Context(), w, Body{
		Code: errorx.OK,
		Msg:  "success",
		Data: data,
	})
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	if codeErr, ok := err.(*errorx.CodeErr); ok {
		httpx.OkJsonCtx(r.Context(), w, Body{
			Code: codeErr.Code,
			Msg:  codeErr.Msg,
		})
		return
	}
	// 非预期的错误，记录详细日志但只返回模糊信息
	logx.Errorf("unexpected error: %v", err)
	httpx.OkJsonCtx(r.Context(), w, Body{
		Code: errorx.ErrUnknown,
		Msg:  "内部错误",
	})
}
