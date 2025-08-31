package repository

import (
	"template-backend/internal/model"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMenuTree(menu model.Menu) ([]*model.Menu, error) {
	var menus []*model.Menu
	db := r.db.Model(&model.Menu{})
	if menu.Type != 0 {
		db = db.Where("type = ?", menu.Type)
	}

	if menu.Visible != nil {
		db = db.Where("visible = ?", *menu.Visible)
	}

	if menu.Name != "" {
		db = db.Where("name LIKE ?", "%"+menu.Name+"%")
	}

	if err := db.Find(&menus).Error; err != nil {
		return nil, err
	}

	// 构建树
	menuMap := make(map[uint]*model.Menu)
	var roots []*model.Menu

	for _, m := range menus {
		menuMap[m.ID] = m
	}

	for _, m := range menus {
		if m.ParentID != nil {
			if parent, ok := menuMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, m)
			}
		} else {
			roots = append(roots, m)
		}
	}
	return roots, nil
}

func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

func (r *MenuRepository) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}
