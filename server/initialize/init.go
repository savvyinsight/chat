package initialize

import (
	"fmt"
	"log"
	"time"

	"chat/global"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(viper.Get("config"))
	fmt.Println(viper.Get("mysql"))
}

func InitMysql() {
	// customize gorm logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, //Slow SQL Threshold
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	global.GVA_DB, _ = gorm.Open(mysql.Open(viper.GetString("Mysql.dns")), &gorm.Config{
		Logger: newLogger,
	})
	// user := model.UserBasic{}
	// global.GVA_DB.Find(&user)
	// fmt.Println(user)
}
