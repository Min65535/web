package model

import (
	"github.com/jinzhu/gorm"
)

// 定义model
type Order struct {
	gorm.Model
	OrderNo  string  `gorm:"unique_index;not null" form:"orderNo" json:"orderNo"`
	UserName string  `form:"username" json:"username"`
	Amount   float64 `form:"amount" json:"amount"`
	Status   string  `form:"status" json:"status"`
	FileUrl  string  `form:"fileUrl" json:"fileUrl"`
}

// 表名字
func (Order) TableName() string {
	return "demo_order"
}
