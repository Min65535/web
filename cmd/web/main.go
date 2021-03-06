package main

import (
	"fmt"
	"web/core/web/dao"
	"web/core/web/handler"
	"web/core/web/router"
	"web/core/web/service"
	"web/pkg/db"
	"web/pkg/setting"
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
