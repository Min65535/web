package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"web/enum"
	"web/model"
	"web/pkg/setting"
	"web/util"
)

var DB *gorm.DB

/**
初始化dbConfig
*/
func InitDB() *gorm.DB {
	var (
		err                                        error
		dbType, dbName, user, password, host, port string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database' :%v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("DATABASE").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	port = sec.Key("PORT").String()
	initConStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", user, password, host, port, "information_schema")
	initCon, err := gorm.Open(dbType, initConStr)
	if err != nil {
		log.Printf("Fail to connect init mysql:%v", err)
		panic("error init connection")
	}

	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"

	if err := initCon.Exec(createDbSQL).Error; err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		panic("error create Db")
	}
	fmt.Println(dbName + "数据库创建成功")
	dbConStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", user, password, host, port, dbName)
	DB, err = gorm.Open(dbType, dbConStr)
	if err != nil {
		log.Printf("Fail to connect mysql:%v", err)
		panic("error db connection")
	}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.AutoMigrate(&model.Order{})

	// 初始化测试数据
	orders := []model.Order{{OrderNo: util.RandomString(10), UserName: "Mike", Amount: 1000.0, Status: enum.Cancel.StatusType, FileUrl: "/home/xxx/picture/cat1.jpeg"},
		{OrderNo: util.RandomString(10), UserName: "Apple", Amount: 10000.0, Status: enum.Doing.StatusType, FileUrl: "/home/xxx/picture/dog1.jpeg"},
		{OrderNo: util.RandomString(10), UserName: "Link", Amount: 6000.0, Status: enum.Finish.StatusType, FileUrl: "/home/xxx/picture/cat4.jpg"},
		{OrderNo: util.RandomString(10), UserName: "Alice", Amount: 18800.0, Status: enum.NotFinish.StatusType, FileUrl: "/home/xxx/picture/cat2.jpeg"},
		{OrderNo: util.RandomString(10), UserName: "Green", Amount: 900.0, Status: enum.NotFinish.StatusType, FileUrl: "/home/xxx/picture/dog3.jpg"}}

	tx := DB.Begin()
	for _, order := range orders {
		result := tx.Create(&order)
		if result.Error != nil {
			tx.Rollback()
			return nil
		}
	}
	tx.Commit()
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
