package model

import (
	"aries/config/db"
	"aries/util"
	"github.com/jinzhu/gorm"
	"strings"
)

// 分类结构
type Category struct {
	gorm.Model
	Children []*Category `gorm:"foreignkey:ParentId" json:"children"`            // 子级分类列表
	ParentId uint        `json:"parent_id"`                                      // 父类 ID：0 表父级分类，大于 0 表子级分类
	Type     uint        `gorm:"type:tinyint(1);unsigned;default:0" json:"type"` // 分类类型，默认值为 0 表文章；1 表友链
	Name     string      `gorm:"type:varchar(100);not null;" json:"name"`        // 分类名称
	Url      string      `gorm:"type:varchar(100)" json:"url"`                   // 访问 URL
}

// 获取所有分类
func (category Category) GetAll(categoryType uint) ([]Category, error) {
	var categories []Category
	var children []Category
	// 查询子分类
	err := db.Db.Where("parent_id > 0").Find(&children).Error
	if err != nil {
		return categories, err
	}
	// 查询所有分类
	err = db.Db.Where("type = ?", categoryType).Find(&categories).Error
	// 将子分类并入父分类
	for i := range categories {
		if categories[i].ParentId == 0 {
			for j := range children {
				if children[j].ParentId == categories[i].ID {
					categories[i].Children = append(categories[i].Children, &children[j])
				}
			}
		}
	}
	return categories, err
}

// 获取分类数据（分页 +　搜索）
func (category Category) GetByPage(page *util.Pagination, key string, categoryType uint) ([]Category, uint, error) {
	var list []Category // 保存结果集
	var children []Category
	// 查询子分类
	err := db.Db.Where("parent_id > 0").Find(&children).Error
	if err != nil {
		return list, 0, err
	}
	// 创建语句
	query := db.Db.Model(&Category{}).Where("type = ?", categoryType).
		Order("created_at desc", true)
	// 拼接搜索语句
	if key != "" {
		query = query.Where("name like concat('%',?,'%')", key)
	}
	// 分页
	total, err := util.ToPage(page, query, &list)
	// 将子类并入父类
	for i := range list {
		if list[i].ParentId == 0 {
			for j := range children {
				if children[j].ParentId == list[i].ID {
					list[i].Children = append(list[i].Children, &children[j])
				}
			}
		}
	}
	// 返回分页数据
	return list, total, err
}

// 获取所有父类
func (category Category) GetAllParents(categoryType uint) (list []Category, err error) {
	err = db.Db.Where("parent_id = 0 and type = ?", categoryType).Find(&list).Error
	return
}

// 添加分类
func (category Category) Create() error {
	return db.Db.Create(&category).Error
}

// 修改分类
func (category Category) Update() error {
	return db.Db.Model(&category).Where("id = ?", category.ID).
		Updates(map[string]interface{}{ // 使用 map 来更新。避免 gorm 默认不更新值 nil, false, 0 的字段
			"name":      category.Name,
			"parent_id": category.ParentId,
			"type":      category.Type,
			"url":       category.Url,
		}).Error
}

// 删除分类
func (category Category) DeleteById(id uint) error {
	// 更新属于该分类的文章的分类 ID 为 null
	err := db.Db.Model(&Article{}).Where("category_id = ?", id).
		Update("category_id", nil).Error
	if err != nil {
		return err
	}
	return db.Db.Where("id = ?", id).
		Unscoped().Delete(&category).Error
}

// 批量删除分类
func (category Category) MultiDelByIds(ids string) error {
	idList := strings.Split(ids, ",") // 根据 , 分割成字符串数组
	// 更新属于该分类的文章的分类 ID 为 null
	err := db.Db.Model(&Article{}).Where("category_id in (?)", idList).
		Update("category_id", nil).Error
	if err != nil {
		return err
	}
	return db.Db.Where("id in (?)", idList).
		Unscoped().Delete(&category).Error
}
