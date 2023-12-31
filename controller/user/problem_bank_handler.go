package user

import (
	"FanCode/controller/utils"
	r "FanCode/models/vo"
	"FanCode/service"
	"github.com/gin-gonic/gin"
)

type ProblemBankController interface {
	// GetAllProblemBank( 读取题库列表
	GetAllProblemBank(ctx *gin.Context)
	// GetProblemBankByID 读取题库信息
	GetProblemBankByID(ctx *gin.Context)
}
type problemBankController struct {
	problemBankService service.ProblemBankService
}

func NewProblemBankController(bankService service.ProblemBankService) ProblemBankController {
	return &problemBankController{
		problemBankService: bankService,
	}
}

func (p *problemBankController) GetAllProblemBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	banks, err := p.problemBankService.GetAllProblemBank()
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(banks)
}

func (p *problemBankController) GetProblemBankByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	bank, err := p.problemBankService.GetProblemBankByID(uint(id))
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(bank)
}
