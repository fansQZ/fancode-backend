package user

import (
	e "FanCode/error"
	"FanCode/models/po"
	r "FanCode/models/vo"
	"FanCode/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProblemController interface {
	// GetProblemList 读取题目列表
	GetProblemList(ctx *gin.Context)
	// GetProblem 读取题目详细信息
	GetProblem(ctx *gin.Context)
	// GetProblemTemplateCode 读取题目编程文件
	GetProblemTemplateCode(ctx *gin.Context)
}

type problemController struct {
	problemService service.ProblemService
}

func NewProblemController(problemService service.ProblemService) ProblemController {
	return &problemController{
		problemService: problemService,
	}
}

func (p *problemController) GetProblemList(ctx *gin.Context) {
	result := r.NewResult(ctx)
	pageQuery, err := GetPageQueryByQuery(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	bankIDStr := ctx.Query("bankID")
	if bankIDStr != "" {
		bankID, err := strconv.Atoi(bankIDStr)
		if err != nil {
			result.Error(e.ErrBadRequest)
			return
		}
		uintBankID := uint(bankID)
		pageQuery.Query = &po.Problem{
			BankID: &uintBankID,
		}
	}
	pageInfo, err := p.problemService.GetUserProblemList(ctx, pageQuery)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(pageInfo)
}

func (p *problemController) GetProblem(ctx *gin.Context) {
	result := r.NewResult(ctx)
	numberStr := ctx.Param("number")
	problem, err := p.problemService.GetProblemByNumber(numberStr)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(problem)
}

func (p *problemController) GetProblemTemplateCode(ctx *gin.Context) {
	result := r.NewResult(ctx)
	problemIDStr := ctx.Param("problemID")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	language := ctx.Param("language")
	codeType := ctx.Param("codeType")
	code, err2 := p.problemService.GetProblemTemplateCode(uint(problemID), language, codeType)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(code)
}