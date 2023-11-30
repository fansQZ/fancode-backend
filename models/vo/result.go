// Package result
// @Author: fzw
// @Create: 2023/7/14
// @Description: 用做统一结果返回的
package vo

import (
	e "FanCode/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

// result
// @Description: 统一result，并返回数据给前端
type result struct {
	ctx *gin.Context
}

func NewResult(ctx *gin.Context) *result {
	return &result{ctx: ctx}
}

// Success1
//
//	@Description: 返回成功结果
//	@receiver r
//	@param data
func (r *result) SuccessData(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := &ResultCont{
		Code:    200,
		Message: "request success",
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *result) SuccessMessage(message string) {
	res := &ResultCont{
		Code:    200,
		Message: message,
		Data:    nil,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *result) Success(message string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := &ResultCont{
		Code:    200,
		Message: message,
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

// 返回异常信息
func (r *result) Error(e *e.Error) {
	res := &ResultCont{
		Code:    e.Code,
		Message: e.Message,
	}
	r.ctx.JSON(e.HttpCode, res)
}

func (r *result) SimpleError(code int, message string, data interface{}) {
	res := &ResultCont{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.ctx.JSON(http.StatusOK, res)
}

func (r *result) SimpleErrorMessage(message string) {
	res := &ResultCont{
		Code:    400,
		Message: message,
	}
	r.ctx.JSON(http.StatusOK, res)
}
