package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"log"
	"mime/multipart"
	"strconv"
	"web/dao"
	"web/db"
	"web/model"
	"web/myerr"
	"web/util"
)

type OrderService struct {
	orderDao *dao.OrderDao
}

func GetOrderService(od *dao.OrderDao) *OrderService {
	return &OrderService{
		orderDao: od,
	}
}

/**
创建order
*/
func (o *OrderService) CreateOrder(order *model.Order) error {
	err := o.orderDao.CreateOrder(order)
	if err != nil {
		log.Printf("Fail to create order because:%v", err)
	}
	return err
}

/**
通过orderno更新username,amount、status、file_url
return 是否更新成功，错误信息
*/
func (o *OrderService) UpdateOrder(order *model.Order) error {
	if len(order.OrderNo) == 0 {
		log.Printf("Fail to UpdateOrder because:orderId is null")
		return myerr.MyError{ErrMsg: "orderId is null"}
	}
	err := o.orderDao.UpdateOrder(order)
	if err != nil {
		log.Printf("Fail to UpdateOrder because:%v", err)
	}
	return err
}

/**
通过orderNo查找Order
*/
func (o *OrderService) GetOrderByOrderNo(orderNo string) (order model.Order, err error) {
	order, err = o.orderDao.GetOrderByOrderNo(orderNo)
	if err != nil {
		log.Printf("Error in OrderService.GetOrderByOrderNo:%v", err)
	}
	return
}

/**
根据user_name模糊查询，order_no，amount，精确查询
根据create_at，amount排序
*/
func (o *OrderService) GetOrderListByColumn(order *model.Order) (orders []model.Order, err error) {
	if len(order.UserName) > 0 {
		orders, err = o.orderDao.GetOrderByOrderNoOrAmountOrUserName(order)
		return
	}
	orders, err = o.orderDao.GetOrderByOrderNoOrAmount(order)
	if err != nil {
		log.Printf("Error in OrderService.GetOrderByOrderNo:%v", err)
	}
	return
}

/**
上传文件
*/
func (o *OrderService) UploadOrderFile(orderNo string, file *multipart.FileHeader, c *gin.Context) (string, error) {
	// 查找orderno是否存在
	_, err := o.GetOrderByOrderNo(orderNo)
	if err != nil {
		log.Printf("Error in OrderService.UploadOrderFile:%v", err)
		return "order不存在", err
	}
	// 开始事务
	tx := db.GetDB().Begin()
	// 获取文家地址
	dst := fmt.Sprintf("/home/hatchersu/picture/%s", util.RandomString(10))
	err = o.orderDao.SaveOrderFileUrl(orderNo, dst, tx)
	if err != nil {
		log.Printf("Error in OrderService.UploadOrderFile:%v", err)
		tx.Rollback() // 保存url出错，事务回滚
	}
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		log.Printf("Error in OrderService.UploadOrderFile:%v", err)
		tx.Rollback() // 文件出错，事务回滚
	}
	tx.Commit() // 提交事务
	// 返回 nil
	return "", err
}

/**
将order导出到xlsx表中
*/
func (o *OrderService) ToExecl() (error, string) {
	outFile := "/home/hatchersu/xlsx/order.xlsx"
	orders, err := o.GetAllOrderList()
	if err != nil {
		return err, ""
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("order")
	if err != nil {
		log.Printf("Error in OrderService.ToExecl:%v", err)
		return err, ""
	}
	initrow := sheet.AddRow()
	IdCell := initrow.AddCell()
	IdCell.Value = "ID"
	CreatedAtCell := initrow.AddCell()
	CreatedAtCell.Value = "CreatedAt"
	UpdatedAtCell := initrow.AddCell()
	UpdatedAtCell.Value = "UpdatedAt"
	DeletedAtCell := initrow.AddCell()
	DeletedAtCell.Value = "DeletedAt"
	OrderNoCell := initrow.AddCell()
	OrderNoCell.Value = "OrderNo"
	UserNameCell := initrow.AddCell()
	UserNameCell.Value = "UserName"
	AmountCell := initrow.AddCell()
	AmountCell.Value = "Amount"
	StatusCell := initrow.AddCell()
	StatusCell.Value = "Status"
	FileUrlCell := initrow.AddCell()
	FileUrlCell.Value = "FileUrl"

	for _, order := range orders {
		row := sheet.AddRow()
		IdCell = row.AddCell()
		IdCell.Value = strconv.Itoa(int(order.ID))
		CreatedAtCell = row.AddCell()
		CreatedAtCell.Value = order.CreatedAt.Format("2006-01-02 15:04:05")
		UpdatedAtCell = row.AddCell()
		UpdatedAtCell.Value = order.UpdatedAt.Format("2006-01-02 15:04:05")
		DeletedAtCell = row.AddCell()
		if order.DeletedAt != nil {
			DeletedAtCell.Value = order.DeletedAt.Format("2006-01-02 15:04:05")
		}
		OrderNoCell = row.AddCell()
		OrderNoCell.Value = order.OrderNo
		UserNameCell = row.AddCell()
		UserNameCell.Value = order.UserName
		AmountCell = row.AddCell()
		AmountCell.Value = fmt.Sprintf("%f", order.Amount)
		StatusCell = row.AddCell()
		StatusCell.Value = order.Status
		FileUrlCell = row.AddCell()
		FileUrlCell.Value = order.FileUrl
	}
	err = file.Save(outFile)
	if err != nil {
		log.Printf("Error in OrderService.ToExecl:%v", err)
	}
	return err, outFile
}

/**
获取所有的orders
*/
func (o *OrderService) GetAllOrderList() (orders []*model.Order, err error) {
	orders, err = o.orderDao.GetAllOrderList()
	if err != nil {
		log.Printf("Error in OrderService.GetAllOrderList:%v", err)
	}
	return
}
