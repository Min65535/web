package main

import (
	"fmt"
	"web/dao"
	"web/db"
	"web/handler"
	"web/pkg/setting"
	"web/router"
	"web/service"
)

func main() {
	db.InitDB()
	DB := db.GetDB()
	defer DB.Close() // 关闭
	// 依赖注入？
	orderDao := dao.GetOrderDao(DB)
	orderService := service.GetOrderService(orderDao)
	orderHandler := handler.GetOrderHandler(orderService)
	r := router.InitRouter(orderHandler)
	r.Run(fmt.Sprintf(":%d", setting.HttpPort))
}
