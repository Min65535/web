package router

import (
	"github.com/gin-gonic/gin"
	"web/handler"
	"web/pkg/setting"
)

/**
初始化路由
*/
func InitRouter(handler *handler.OrderHandler) *gin.Engine {
	router := gin.Default()
	gin.SetMode(setting.RunMode)
	api := router.Group("/order")
	{
		// 获取order通过orderno
		api.GET("/get_by_order_no", handler.GetOrder)
		// 获取orderlist通过username、等
		api.GET("/get_order_list", handler.GetOrderList)
		// 更新order
		api.PUT("/update_order", handler.UpdateOrder)
		// 创建order
		api.POST("/save_order", handler.CreateOrder)
		// 上传文件
		api.POST("/upload_file", handler.UploadFile)
		// 下载文件
		api.GET("/download_file", handler.FileDownload)
		// 导出excel
		api.GET("/to_excel", handler.OrderToExcel)
	}
	return router
}
