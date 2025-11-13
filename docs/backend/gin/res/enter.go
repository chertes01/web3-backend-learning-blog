package res

import "github.com/gin-gonic/gin"

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

var codeMap = map[int]string{
	1001: "参数错误",
	1002: "未授权",
	1003: "禁止访问",
	1004: "资源未找到",
	1005: "服务错误",
}

func response(c *gin.Context, code int, data any, msg string) {
	c.JSON(200, Response{
		Msg:  msg,
		Code: code,
		Data: data,
	})
}

func Ok(c *gin.Context, data any, msg string) {
	response(c, 0, data, msg)
}

func OkData(c *gin.Context, data any) {
	Ok(c, data, "success")
}

func OkMsg(c *gin.Context, msg string) {
	Ok(c, gin.H{}, msg)
}

func Fail(c *gin.Context, code int, data any, msg string) {
	response(c, code, data, msg)
}

func FailMsg(c *gin.Context, msg string) {
	Fail(c, 1001, nil, msg) // 默认错误码1001
}

func FailCode(c *gin.Context, code int, msg string) {
	msg, ok := codeMap[code]
	if !ok {
		msg = "未知错误"
	}
	Fail(c, code, nil, msg) //传入空数据
}

func FailData(c *gin.Context, data any) {
	Fail(c, 1001, data, "")
}
