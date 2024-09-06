package main

import (
	router "ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	r := router.Router()
	r.Run(":8081")
}
