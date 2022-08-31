/**
    @auther: oreki
    @date: 2022/6/20
    @note: 图灵老祖保佑,永无BUG
**/

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"srv/user_srv/model"
	"time"
)

var (
	DB *gorm.DB
)

func init() {

	dsn := "root:root@tcp(42.192.220.243:3306)/gp_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	_ = DB.AutoMigrate(&model.User{}) //此处应该有sql语句
}
