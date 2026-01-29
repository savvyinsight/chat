package main

import (
	"chat/initialize"
	"chat/router"
)

func main() {
	initialize.InitConfig()
	initialize.InitMysql()
	initialize.InitRedis()
	r := router.Router()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
