package service

import (
	"FanCode/dao/mock"
	e "FanCode/error"
	"FanCode/global"
	"FanCode/models/dto"
	"FanCode/models/po"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

func TestSysApiService_GetApiCount(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	apiDao := mock.NewMockSysApiDao(mockCtl)
	apiDao.EXPECT().GetApiCount(gomock.Any()).Return(int64(10), nil)
	apiDao.EXPECT().GetApiCount(gomock.Any()).Return(int64(0), gorm.ErrInvalidDB)
	apiService := NewSysApiService(apiDao)
	count, err := apiService.GetApiCount()
	assert.Equal(t, int64(10), count)
	assert.Nil(t, err)
	count, err = apiService.GetApiCount()
	assert.Equal(t, err, e.ErrMysql)
	assert.Equal(t, int64(0), count)

}

func TestSysApiService_DeleteApiByID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	apiDao := mock.NewMockSysApiDao(mockCtl)

	// mock数据库
	db, mock, err := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectCommit()
	assert.Nil(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	assert.Nil(t, err)
	global.Mysql = gormDB

	// 生成api树
	apis1 := make([]*po.SysApi, 8)
	apis2 := make([]*po.SysApi, 8)
	for i := 0; i < 8; i++ {
		api := &po.SysApi{}
		api.ID = uint(i + 2)
		apis1[i] = api
		api = &po.SysApi{}
		api.ID = uint(i + 10)
		apis2[i] = api
	}
	apiDao.EXPECT().GetChildApisByParentID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db2 *gorm.DB, id uint) ([]*po.SysApi, *e.Error) {
			if id == 1 {
				return apis1, nil
			}
			if id == 9 {
				return apis2, nil
			}
			return []*po.SysApi{}, nil
		}).AnyTimes()

	deleteCount := 0
	apiDao.EXPECT().DeleteApiByID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(db *gorm.DB, apiID uint) error {
			deleteCount += int(apiID)
			return nil
		}).AnyTimes()

	apiService := NewSysApiService(apiDao)
	err2 := apiService.DeleteApiByID(1)
	assert.Nil(t, err2)
	assert.Equal(t, deleteCount, ((1+17)/2)*17)
}

func TestSysApiService_GetApiByID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	api.ID = 1
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(1)).Return(api, nil)
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(3)).Return(nil, gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	api2, err := apiService.GetApiByID(1)
	assert.Equal(t, api, api2)
	assert.Nil(t, err)

	api3, err := apiService.GetApiByID(2)
	assert.Nil(t, api3)
	assert.Equal(t, err, e.ErrApiNotExist)

	api4, err := apiService.GetApiByID(3)
	assert.Nil(t, api4)
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_UpdateApi(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	api.ID = 1
	apiDao.EXPECT().UpdateApi(gomock.Any(), api).
		DoAndReturn(func(db *gorm.DB, sysApi *po.SysApi) error {
			assert.NotEqual(t, sysApi.UpdatedAt, time.Time{})
			sysApi.UpdatedAt = time.Time{}
			assert.Equal(t, api, sysApi)
			return nil
		})
	apiDao.EXPECT().UpdateApi(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	err := apiService.UpdateApi(api)
	assert.Nil(t, err)
	err = apiService.UpdateApi(&po.SysApi{})
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_InsertApi(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	apiDao.EXPECT().InsertApi(gomock.Any(), api).
		DoAndReturn(func(db *gorm.DB, sysApi *po.SysApi) error {
			assert.Equal(t, api, sysApi)
			sysApi.ID = 1
			return nil
		})
	apiDao.EXPECT().InsertApi(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	id, err := apiService.InsertApi(api)
	assert.Equal(t, id, uint(1))
	assert.Nil(t, err)
	id, err = apiService.InsertApi(&po.SysApi{})
	assert.Equal(t, id, uint(0))
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_GetApiTree(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	apis := make([]*po.SysApi, 4)
	for i := 0; i < 4; i++ {
		api := &po.SysApi{}
		api.Name = "api" + strconv.Itoa(i)
		api.Description = "api描述" + strconv.Itoa(i)
		api.Path = "apiPath" + strconv.Itoa(i)
		api.Method = "apiMethod" + strconv.Itoa(i)
		api.ID = uint(i + 1)
		apis[i] = api
	}
	apis[1].ParentApiID = 1
	apis[2].ParentApiID = 1
	apis[3].ParentApiID = 1
	apiDao.EXPECT().GetAllApi(gomock.Any()).Return(apis, nil)
	apiDao.EXPECT().GetAllApi(gomock.Any()).Return([]*po.SysApi{}, gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)

	treeDtos, err := apiService.GetApiTree()
	treeDto := dto.NewSysApiTreeDto(apis[0])
	for i := 1; i < 4; i++ {
		treeDto.Children = append(treeDto.Children, dto.NewSysApiTreeDto(apis[i]))
	}
	assert.Equal(t, []*dto.SysApiTreeDto{treeDto}, treeDtos)
	assert.Nil(t, err)

	treeDtos, err = apiService.GetApiTree()
	assert.Nil(t, treeDtos)
	assert.Equal(t, err, e.ErrMysql)
}
