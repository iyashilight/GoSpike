package errorx

const (
	OK                     = 0
	ErrUserNotFound        = 1001
	ErrInvalidCredential   = 1002
	ErrUserAlreadyExists   = 1003
	ErrProductNotFound     = 2001
	ErrStockNotEnough      = 2002
	ErrCreateProductFailed = 2003
	ErrUpdateProductFailed = 2004
	ErrDeleteProductFailed = 2005
	ErrCacheStockFailed    = 2006
	ErrOrderNotFound       = 3001
	ErrCreateOrderFailed   = 3002
	ErrUnknown             = 9999
)

type CodeErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeErr) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg ...string) *CodeErr {
	m := map[int]string{
		OK:                     "success",
		ErrUserNotFound:        "用户不存在",
		ErrInvalidCredential:   "用户名或密码错误",
		ErrUserAlreadyExists:   "用户已存在",
		ErrProductNotFound:     "商品不存在",
		ErrStockNotEnough:      "库存不足",
		ErrCreateProductFailed: "创建商品失败",
		ErrUpdateProductFailed: "更新商品失败",
		ErrDeleteProductFailed: "删除商品失败",
		ErrCacheStockFailed:    "缓存库存失败",
		ErrOrderNotFound:       "订单不存在",
		ErrCreateOrderFailed:   "创建订单失败",
		ErrUnknown:             "未知错误",
	}
	if len(msg) > 0 {
		return &CodeErr{Code: code, Msg: msg[0]}
	}
	v, ok := m[code]
	if !ok {
		return &CodeErr{Code: ErrUnknown, Msg: m[ErrUnknown]}
	}
	return &CodeErr{Code: code, Msg: v}
}
