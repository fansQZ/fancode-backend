package admin

import (
	"FanCode/controller/utils"
	e "FanCode/error"
	"FanCode/models/po"
	r "FanCode/models/vo"
	"FanCode/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ProblemBankManagementController
// @Description: 题库管理相关功能
type ProblemBankManagementController interface {
	// UploadProblemBankIcon 上传题库图标
	UploadProblemBankIcon(ctx *gin.Context)
	// ReadProblemBankIcon 读取题库图标
	ReadProblemBankIcon(ctx *gin.Context)
	// InsertProblemBank 添加题库
	InsertProblemBank(ctx *gin.Context)
	// UpdateProblemBank 更新题库
	UpdateProblemBank(ctx *gin.Context)
	// DeleteProblemBank 删除题库
	DeleteProblemBank(ctx *gin.Context)
	// GetProblemBankList 读取题库列表
	GetProblemBankList(ctx *gin.Context)
	// GetSimpleProblemBankList 读取简单的题库列表
	GetSimpleProblemBankList(ctx *gin.Context)
	// GetProblemBankByID 读取题库信息
	GetProblemBankByID(ctx *gin.Context)
}

type problemBankManagementController struct {
	problemBankService service.ProblemBankService
}

func NewProblemBankManagementController(bankService service.ProblemBankService) ProblemBankManagementController {
	return &problemBankManagementController{
		problemBankService: bankService,
	}
}

func (p *problemBankManagementController) UploadProblemBankIcon(ctx *gin.Context) {
	result := r.NewResult(ctx)
	file, err := ctx.FormFile("icon")
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if file.Size > 2<<20 {
		result.SimpleErrorMessage("文件大小不能超过2m")
		return
	}
	path, err2 := p.problemBankService.UploadProblemBankIcon(file)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(path)
}

func (p *problemBankManagementController) ReadProblemBankIcon(ctx *gin.Context) {
	avatarName := ctx.Param("iconName")
	p.problemBankService.ReadProblemBankIcon(ctx, avatarName)
}

func (p *problemBankManagementController) InsertProblemBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	bank := p.getBank(ctx)
	pID, err := p.problemBankService.InsertProblemBank(bank, ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.Success("题库添加成功", pID)
}

func (p *problemBankManagementController) UpdateProblemBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	bank := p.getBank(ctx)
	// 读取id
	bankIDString := ctx.PostForm("id")
	bankID, err := strconv.Atoi(bankIDString)
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	bank.ID = uint(bankID)
	// 读取文件
	err2 := p.problemBankService.UpdateProblemBank(bank)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData("题库修改成功")
}

func (p *problemBankManagementController) DeleteProblemBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	// 读取id
	bankID := uint(utils.GetIntParamOrDefault(ctx, "id", 0))
	// 判断是否强制删除
	forceDeleteStr := ctx.Param("forceDelete")
	forceDelete := false
	if forceDeleteStr == "true" {
		forceDelete = true
	}
	// 删除题库
	err2 := p.problemBankService.DeleteProblemBank(bankID, forceDelete)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData("题库删除成功")
}

func (p *problemBankManagementController) GetProblemBankList(ctx *gin.Context) {
	result := r.NewResult(ctx)
	pageQuery, err := utils.GetPageQueryByQuery(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	// 读取名称和描述
	bank := &po.ProblemBank{
		Name:        ctx.Query("name"),
		Description: ctx.Query("description"),
	}
	pageQuery.Query = bank
	pageInfo, err := p.problemBankService.GetProblemBankList(pageQuery)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(pageInfo)
}

func (p *problemBankManagementController) GetSimpleProblemBankList(ctx *gin.Context) {
	result := r.NewResult(ctx)
	banks, err := p.problemBankService.GetSimpleProblemBankList()
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(banks)
}

func (p *problemBankManagementController) GetProblemBankByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	bank, err2 := p.problemBankService.GetProblemBankByID(uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(bank)
}

func (p *problemBankManagementController) getBank(ctx *gin.Context) *po.ProblemBank {
	bank := &po.ProblemBank{}
	bank.Name = ctx.PostForm("name")
	bank.Description = ctx.PostForm("description")
	bank.Icon = ctx.PostForm("icon")
	return bank
}
