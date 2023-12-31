package controller

import (
	e "FanCode/error"
	"FanCode/models/dto"
	"FanCode/models/po"
	r "FanCode/models/vo"
	"FanCode/service"
	"FanCode/utils"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	// Login 用户登录
	Login(ctx *gin.Context)
	// SendAuthCode 发送验证码
	SendAuthCode(ctx *gin.Context)
	// UserRegister 用户注册
	UserRegister(ctx *gin.Context)
	// GetUserInfo 从token里面读取用户信息
	GetUserInfo(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (u *authController) SendAuthCode(ctx *gin.Context) {
	result := r.NewResult(ctx)
	email := ctx.PostForm("email")
	kind := ctx.PostForm("type")
	if email != "" && !utils.VerifyEmailFormat(email) {
		result.SimpleErrorMessage("邮箱格式错误")
		return
	}
	// 生成code
	if _, err := u.authService.SendAuthCode(email, kind); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("验证码发送成功")
}

func (u *authController) UserRegister(ctx *gin.Context) {
	result := r.NewResult(ctx)
	user := &po.SysUser{
		Email:    ctx.PostForm("email"),
		Username: ctx.PostForm("username"),
		Password: ctx.PostForm("password"),
	}
	code := ctx.PostForm("code")
	if err := u.authService.UserRegister(user, code); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("注册成功")
}

func (u *authController) Login(ctx *gin.Context) {
	result := r.NewResult(ctx)
	//获取并检验用户参数
	kind := ctx.PostForm("type")
	account := ctx.PostForm("account")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	code := ctx.PostForm("code")
	if kind != "password" && kind != "email" {
		result.Error(e.ErrBadRequest)
		return
	} else if kind == "password" && (account == "" || password == "") {
		result.Error(e.ErrBadRequest)
		return
	} else if kind == "email" && (email == "" || code == "") {
		result.Error(e.ErrBadRequest)
		return
	}
	// 登录
	var token string
	var err *e.Error
	if kind == "password" {
		token, err = u.authService.PasswordLogin(account, password)
	} else if kind == "email" {
		token, err = u.authService.EmailLogin(email, code)
	} else {
		result.Error(e.ErrLoginType)
		return
	}
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(token)

}

func (u *authController) GetUserInfo(ctx *gin.Context) {
	result := r.NewResult(ctx)
	user := ctx.Keys["user"].(*dto.UserInfo)
	result.SuccessData(user)
}
