package initialize

import (
	"fmt"
	"log"

	"chat/global"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	log.New(log.Writer(), "\r\n", log.LstdFlags) // io writer

	global.GVA_DB, _ = gorm.Open(mysql.Open(viper.GetString("Mysql.dns")), &gorm.Config{})
	// user := model.UserBasic{}
	// global.GVA_DB.Find(&user)
	// fmt.Println(user)
}
