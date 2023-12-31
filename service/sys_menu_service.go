package service

import (
	"FanCode/dao"
	e "FanCode/error"
	"FanCode/global"
	"FanCode/models/dto"
	"FanCode/models/po"
	"gorm.io/gorm"
	"time"
)

type SysMenuService interface {

	// GetMenuCount 获取menu数目
	GetMenuCount() (int64, *e.Error)
	// DeleteMenuByID 删除menu
	DeleteMenuByID(id uint) *e.Error
	// UpdateMenu 更新menu
	UpdateMenu(menu *po.SysMenu) *e.Error
	// GetMenuByID 根据id获取menu
	GetMenuByID(id uint) (*po.SysMenu, *e.Error)
	// GetMenuTree 获取menu树
	GetMenuTree() ([]*dto.SysMenuTreeDto, *e.Error)
	// InsertMenu 添加menu
	InsertMenu(menu *po.SysMenu) (uint, *e.Error)
}

type sysMenuService struct {
	sysMenuDao dao.SysMenuDao
}

func NewSysMenuService(menuDao dao.SysMenuDao) SysMenuService {
	return &sysMenuService{
		sysMenuDao: menuDao,
	}
}

func (s *sysMenuService) GetMenuCount() (int64, *e.Error) {
	count, err := s.sysMenuDao.GetMenuCount(global.Mysql)
	if err != nil {
		return 0, e.ErrMysql
	}
	return count, nil
}

// DeleteMenuByID 根据menu的id进行删除
func (s *sysMenuService) DeleteMenuByID(id uint) *e.Error {
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 递归删除API
		return s.deleteMenusRecursive(tx, id)
	})

	if err != nil {
		return e.ErrMysql
	}

	return nil
}

// deleteMenusRecursive 递归删除API
func (s *sysMenuService) deleteMenusRecursive(db *gorm.DB, parentID uint) error {
	childMenus, err := s.sysMenuDao.GetChildMenusByParentID(db, parentID)
	if err != nil {
		return err
	}
	for _, childAPI := range childMenus {
		// 删除子menu的子menu
		if err = s.deleteMenusRecursive(db, childAPI.ID); err != nil {
			return err
		}
	}
	// 当前menu
	if err = s.sysMenuDao.DeleteMenuByID(db, parentID); err != nil {
		return err
	}
	return nil
}

func (s *sysMenuService) UpdateMenu(menu *po.SysMenu) *e.Error {
	menu.UpdatedAt = time.Now()
	err := s.sysMenuDao.UpdateMenu(global.Mysql, menu)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *sysMenuService) GetMenuByID(id uint) (*po.SysMenu, *e.Error) {
	menu, err := s.sysMenuDao.GetMenuByID(global.Mysql, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, e.ErrMysql
	}
	if err == gorm.ErrRecordNotFound {
		return nil, e.ErrMenuNotExist
	}
	return menu, nil
}

func (s *sysMenuService) GetMenuTree() ([]*dto.SysMenuTreeDto, *e.Error) {
	var menuList []*po.SysMenu
	var err error
	if menuList, err = s.sysMenuDao.GetAllMenu(global.Mysql); err != nil {
		return nil, e.ErrMysql
	}

	menuMap := make(map[uint]*dto.SysMenuTreeDto)
	var rootMenus []*dto.SysMenuTreeDto

	// 添加到map中保存
	for _, menu := range menuList {
		menuMap[menu.ID] = dto.NewSysMenuTreeDto(menu)
	}

	// 遍历并添加到父节点中
	for _, menu := range menuList {
		if menu.ParentMenuID == 0 {
			rootMenus = append(rootMenus, menuMap[menu.ID])
		} else {
			parentMenu, exists := menuMap[menu.ParentMenuID]
			if !exists {
				return nil, e.ErrMenuUnknownError
			}
			parentMenu.Children = append(parentMenu.Children, menuMap[menu.ID])
		}
	}

	return rootMenus, nil
}

func (s *sysMenuService) InsertMenu(menu *po.SysMenu) (uint, *e.Error) {
	err := s.sysMenuDao.InsertMenu(global.Mysql, menu)
	if err != nil {
		return 0, e.ErrMysql
	}
	return menu.ID, nil
}
