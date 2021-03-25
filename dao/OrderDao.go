package dao

import (
	"github.com/jinzhu/gorm"
	"web/model"
	"web/util"
)

type OrderDao struct {
	DB *gorm.DB
}

func GetOrderDao(DB *gorm.DB) *OrderDao{
	return &OrderDao{
		DB: DB,
	}
}

func (od *OrderDao) CreateOrder(order *model.Order) error {
	order.OrderNo = util.RandomString(15)
	err := od.DB.Create(&order).Error
	return err
}

func  (od *OrderDao) UpdateOrder(order *model.Order) error{
	err := od.DB.Model(&model.Order{}).Where("order_no = ?",order.OrderNo).Select("UserName","Amount","Status","FileUrl").Updates(order).Error
	return err
}

func (od *OrderDao) GetOrderByOrderNo(orderNo string) (order model.Order,err error)   {
	order = model.Order{}
	err = od.DB.Where("order_no = ?",orderNo).Take(&order).Error
	return
}

func (od *OrderDao) GetOrderByOrderNoOrAmount(order *model.Order)(orders []model.Order,err error) {
	err = od.DB.Where("order_no = ? ",order.OrderNo).
		Or("amount = ? ",order.Amount).
		Order("created_at , amount").
		Find(&orders).Error
	return
}

func (od *OrderDao) GetOrderByOrderNoOrAmountOrUserName(order *model.Order)(orders []model.Order,err error) {
	err = od.DB.Where("user_name like ? ","%"+order.UserName+"%").
		Or("order_no = ? ",order.OrderNo).
		Or("amount = ? ",order.Amount).
		Order("created_at , amount").
		Find(&orders).Error
	return
}

func  (od *OrderDao) GetAllOrderList() (orders []*model.Order,err error){
	result := od.DB.Find(&orders)
	err = result.Error
	return
}
/**
 保存fileUrl
 */
func (od *OrderDao) SaveOrderFileUrl(orderNo string,fileUrl string,tx *gorm.DB) error{
	err := tx.Model(&model.Order{}).Where("order_no = ?", orderNo).Update(&model.Order{FileUrl: fileUrl}).Error
	return err
}
