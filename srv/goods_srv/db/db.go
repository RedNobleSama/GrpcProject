/**
    @auther: oreki
    @date: 2022/4/27
    @note: 图灵老祖保佑,永无BUG
**/

package db

import (
	"goods_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

func init() {
	dsn := "root:Aa123456@tcp(127.0.0.1:3306)/gp_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢SQL阈值
			LogLevel:      logger.Silent, // 日记等级
			Colorful:      true,          // 禁用彩色打印
		},
	)
	// 全局模式
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表面
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//hashPWS, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	//
	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("cyt%d", i),
	//		Mobile:   fmt.Sprintf("187815644%d", i),
	//		Password: string(hashPWS),
	//	}
	//	DB.Save(&user)
	//}

	//设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	//sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么

	//定义一个表结构， 将表结构直接生成对应的表 - migrations
	// 迁移 schema
	_ = DB.AutoMigrate(&model.Banner{}, &model.Category{}, &model.GoodsCategoryBrand{}, &model.Brands{}, &model.Goods{}) //此处应该有sql语句
}
