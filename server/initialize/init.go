package initialize

import (
	"context"
	"fmt"
	"log"
	"time"

	"chat/global"

	"github.com/go-redis/redis/v8"
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

func InitRedis() {
	addr := viper.GetString("Redis.Addr")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	pwd := viper.GetString("Redis.Password")
	db := viper.GetInt("Redis.DB")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("redis ping failed: %v", err)
	} else {
		log.Printf("redis connected: %s", addr)
	}

	global.GVA_REDIS = client
	// store context as global.GVA_CTX already set in global package
	_ = ctx
}
