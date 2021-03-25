package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"web/core/web/service"
	"web/pkg/enum"
	"web/pkg/model"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func GetOrderHandler(os *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

/**
通过orderno获取order
*/
func (oh *OrderHandler) GetOrder(c *gin.Context) {
	orderNo := c.Query("orderNo")
	if len(orderNo) == 0 {
		c.JSON(http.StatusOK, enum.ParamError.ToFailGinHWithMsg("订单ID的不能为空"))
		return
	}
	order, err := oh.orderService.GetOrderByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusOK, enum.ServerError.ToFailGinHWithMsg("该order不存在"))
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinHWithMsgAndData("", order))
}

/**
通过orderno、username、amount获取orderlist 并且根据amount排序
*/
func (oh *OrderHandler) GetOrderList(c *gin.Context) {
	orderNo := c.Query("orderNo")
	username := c.Query("username")
	amountStr := c.Query("amount")
	var amount float64
	// 判断金额是否为空，不为空则转为浮点数
	if len(amountStr) != 0 {
		var err error
		amount, err = strconv.ParseFloat(amountStr, 64) // 将字符串转化为float64
		if err != nil {
			c.JSON(http.StatusOK, enum.ParamError.ToFailGinHWithMsg("金额必须为浮点数"))
			return
		}
	}
	order := model.Order{
		OrderNo:  orderNo,
		UserName: username,
		Amount:   amount,
	}
	orders, err := oh.orderService.GetOrderListByColumn(&order)
	if err != nil {
		c.JSON(http.StatusOK, enum.TotalError.ToFailGinHWithMsg("查询的不存在"))
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinHWithMsgAndData("", orders))
}

/**
更新order
*/
func (oh *OrderHandler) UpdateOrder(c *gin.Context) {
	var order model.Order
	err := c.ShouldBind(&order)
	if err != nil {
		log.Printf("error about paresFloat in OrderController.UpdateOrder: %v", err.Error())
		c.JSON(http.StatusBadRequest, enum.ParamError.ToFailGinH())
		return
	}
	// 进行更新
	err = oh.orderService.UpdateOrder(&order)
	if err != nil {
		c.JSON(http.StatusOK, enum.NullParam.ToFailGinHWithMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinH())
}

/**
创建order
*/
func (oh *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBind(&order); err != nil {
		log.Printf("error aber.out paresFloat in OrderControllUpdateOrder: %v", err.Error())
		c.JSON(http.StatusOK, enum.ParamError.ToFailGinH())
		return
	}
	// 进行创建
	err := oh.orderService.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusOK, enum.ServerError.ToFailGinH())
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinH())
}

/**
上传文件
*/
func (oh *OrderHandler) UploadFile(c *gin.Context) {
	form, err := c.MultipartForm() // 获取上传的数据
	if err != nil {
		log.Printf("Fail to uploadFile in OrderController.uploadFile:%v", err.Error())
		c.JSON(http.StatusOK, enum.UploadError.ToFailGinH())
		return
	}
	file := form.File["file"][0]
	orderNo := form.Value["orderNo"][0]
	// 判断orderNo是否为空
	if len(orderNo) == 0 {
		c.JSON(http.StatusOK, enum.NullParam.ToFailGinH())
		return
	}
	msg, err := oh.orderService.UploadOrderFile(orderNo, file, c)
	if err != nil {
		c.JSON(http.StatusOK, enum.UploadError.ToFailGinHWithMsg(msg))
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinH())
}

/**
下载文件
*/
func (oh *OrderHandler) FileDownload(c *gin.Context) {
	orderNo := c.Query("orderNo")
	// 查询orderId对应的order
	order, err := oh.orderService.GetOrderByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusOK, enum.TotalError.ToFailGinHWithMsg("orderNo没有对应的订单"))
		return
	}
	// 打开文件路径
	fileTmp, err := os.Open(order.FileUrl)
	if err != nil {
		c.JSON(http.StatusOK, enum.TotalError.ToFailGinHWithMsg("文件打开错误！"))
		log.Printf("Fail in FileDownload: FileOpen Error: %v", err)
		return
	}
	defer fileTmp.Close()
	// 打开文件
	fileName := path.Base(order.FileUrl)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	if len(fileName) == 0 || len(order.FileUrl) == 0 {
		log.Printf("Fail in FileDownload: fileName=%s,FileUrl=%s", fileName, order.FileUrl)
		c.JSON(http.StatusOK, enum.TotalError.ToFailGinHWithMsg("获取文件资源失败"))
		return
	}
	c.File(order.FileUrl)
}

func (oh *OrderHandler) OrderToExcel(c *gin.Context) {
	err, dst := oh.orderService.ToExecl()
	if err != nil {
		c.JSON(http.StatusOK, enum.TotalError.ToFailGinHWithMsg("获取资源失败"))
		return
	}
	c.JSON(http.StatusOK, enum.Success.ToSuccessGinHWithMsgAndData("成功导入到目录："+dst, nil))
}
