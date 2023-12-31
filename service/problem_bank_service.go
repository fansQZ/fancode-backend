package service

import (
	conf "FanCode/config"
	"FanCode/dao"
	e "FanCode/error"
	"FanCode/file_store"
	"FanCode/global"
	"FanCode/models/dto"
	"FanCode/models/po"
	r "FanCode/models/vo"
	"FanCode/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
	"time"
)

const (
	// ProblemBankIconPath cos中，题库图标存储位置
	ProblemBankIconPath = "/icon/problemBank"
)

// ProblemBankService 题库管理的service
type ProblemBankService interface {
	// UploadProblemBankIcon 上传题库图标
	UploadProblemBankIcon(file *multipart.FileHeader) (string, *e.Error)
	// ReadProblemBankIcon 读取题库图标
	ReadProblemBankIcon(ctx *gin.Context, iconName string)
	// InsertProblemBank 添加题库
	InsertProblemBank(problemBank *po.ProblemBank, ctx *gin.Context) (uint, *e.Error)
	// UpdateProblemBank 更新题库
	UpdateProblemBank(problemBank *po.ProblemBank) *e.Error
	// DeleteProblemBank 删除题库
	DeleteProblemBank(id uint, forceDelete bool) *e.Error
	// GetProblemBankList 获取题库列表
	GetProblemBankList(query *dto.PageQuery) (*dto.PageInfo, *e.Error)
	// GetAllProblemBank 获取所有的题库列表
	GetAllProblemBank() ([]*dto.ProblemBankDtoForList, *e.Error)
	// GetSimpleProblemBankList 获取简单的题库列表
	GetSimpleProblemBankList() ([]*dto.ProblemBankDtoForSimpleList, *e.Error)
	// GetProblemBankByID 获取题库信息
	GetProblemBankByID(id uint) (*po.ProblemBank, *e.Error)
}

type problemBankService struct {
	config         *conf.AppConfig
	problemBankDao dao.ProblemBankDao
	problemDao     dao.ProblemDao
	sysUserDao     dao.SysUserDao
}

func NewProblemBankService(config *conf.AppConfig, bankDao dao.ProblemBankDao, problemDao dao.ProblemDao, userDao dao.SysUserDao) ProblemBankService {
	return &problemBankService{
		config:         config,
		problemBankDao: bankDao,
		problemDao:     problemDao,
		sysUserDao:     userDao,
	}
}

func (p *problemBankService) UploadProblemBankIcon(file *multipart.FileHeader) (string, *e.Error) {
	cos := file_store.NewImageCOS(p.config.COSConfig)
	fileName := file.Filename
	fileName = utils.GetUUID() + "." + path.Base(fileName)
	file2, err := file.Open()
	if err != nil {
		return "", e.ErrBadRequest
	}
	err = cos.SaveFile(path.Join(ProblemBankIconPath, fileName), file2)
	if err != nil {
		return "", e.ErrServer
	}
	return p.config.ProUrl + path.Join("/manage/problemBank/icon", fileName), nil
}

func (p *problemBankService) ReadProblemBankIcon(ctx *gin.Context, iconName string) {
	result := r.NewResult(ctx)
	cos := file_store.NewImageCOS(p.config.COSConfig)
	bytes, err := cos.ReadFile(path.Join(ProblemBankIconPath, iconName))
	if err != nil {
		result.Error(e.ErrServer)
		return
	}
	_, _ = ctx.Writer.Write(bytes)
}

func (p *problemBankService) InsertProblemBank(problemBank *po.ProblemBank, ctx *gin.Context) (uint, *e.Error) {
	// 对设置值的数据设置默认值
	if problemBank.Name == "" {
		problemBank.Name = "未命名题库"
	}
	if problemBank.Description == "" {
		problemBank.Description = "无描述信息"
	}
	problemBank.CreatorID = ctx.Keys["user"].(*dto.UserInfo).ID
	err := p.problemBankDao.InsertProblemBank(global.Mysql, problemBank)
	if err != nil {
		return 0, e.ErrMysql
	}
	return problemBank.ID, nil
}

func (p *problemBankService) UpdateProblemBank(problemBank *po.ProblemBank) *e.Error {
	problemBank.CreatorID = 0
	problemBank.UpdatedAt = time.Now()
	err := p.problemBankDao.UpdateProblemBank(global.Mysql, problemBank)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (p *problemBankService) DeleteProblemBank(id uint, forceDelete bool) *e.Error {
	var err error
	// 非强制删除
	if !forceDelete {
		var count int64
		count, err = p.problemDao.GetProblemCount(global.Mysql, &po.Problem{
			BankID: &id,
		})
		if count != 0 {
			return e.NewCustomMsg("题库不为空，请问是否需要强制删除")
		}
		err = p.problemBankDao.DeleteProblemBankByID(global.Mysql, id)
		if err != nil {
			return e.ErrMysql
		}
		return nil
	}

	// 强制删除
	err = p.problemBankDao.DeleteProblemBankByID(global.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (p *problemBankService) GetProblemBankList(query *dto.PageQuery) (*dto.PageInfo, *e.Error) {
	var bankQuery *po.ProblemBank
	if query.Query != nil {
		bankQuery = query.Query.(*po.ProblemBank)
	}
	// 获取题库列表
	banks, err := p.problemBankDao.GetProblemBankList(global.Mysql, query)
	if err != nil {
		return nil, e.ErrMysql
	}
	newProblemBanks := make([]*dto.ProblemBankDtoForList, len(banks))
	for i := 0; i < len(banks); i++ {
		newProblemBanks[i] = dto.NewProblemBankDtoForList(banks[i])
		// 读取题库中的题目总数还有作者
		newProblemBanks[i].ProblemCount, err = p.problemDao.GetProblemCount(global.Mysql, &po.Problem{
			BankID: &newProblemBanks[i].ID,
		})
		if err != nil {
			return nil, e.ErrMysql
		}
		newProblemBanks[i].CreatorName, err = p.sysUserDao.GetUserNameByID(global.Mysql, banks[i].CreatorID)
	}
	// 获取所有题库总数目
	var count int64
	count, err = p.problemBankDao.GetProblemBankCount(global.Mysql, bankQuery)
	if err != nil {
		return nil, e.ErrMysql
	}
	pageInfo := &dto.PageInfo{
		Total: count,
		Size:  int64(len(newProblemBanks)),
		List:  newProblemBanks,
	}
	return pageInfo, nil
}

func (p *problemBankService) GetAllProblemBank() ([]*dto.ProblemBankDtoForList, *e.Error) {
	banks, err := p.problemBankDao.GetAllProblemBank(global.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	answer := make([]*dto.ProblemBankDtoForList, len(banks))
	for index, bank := range banks {
		answer[index] = dto.NewProblemBankDtoForList(bank)
	}
	return answer, nil
}

func (p *problemBankService) GetSimpleProblemBankList() ([]*dto.ProblemBankDtoForSimpleList, *e.Error) {
	banks, err := p.problemBankDao.GetSimpleProblemBankList(global.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	newBanks := make([]*dto.ProblemBankDtoForSimpleList, len(banks))
	for i := 0; i < len(banks); i++ {
		newBanks[i] = dto.NewProblemBankDtoForSimpleList(banks[i])
	}
	return newBanks, nil
}

func (p *problemBankService) GetProblemBankByID(id uint) (*po.ProblemBank, *e.Error) {
	bank, err := p.problemBankDao.GetProblemBankByID(global.Mysql, id)
	if err == gorm.ErrRecordNotFound {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return bank, nil
}
