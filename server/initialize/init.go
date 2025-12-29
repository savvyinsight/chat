package initialize

import (
	"fmt"

	"chat/global"
	"chat/model"

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
	global.GVA_DB, _ = gorm.Open(mysql.Open(viper.GetString("Mysql")), &gorm.Config{})
	user := model.UserBasic{}
	global.GVA_DB.Find(&user)
	fmt.Println(user)
}
