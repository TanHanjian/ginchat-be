package main

import (
	router "ginchat/router"
	"ginchat/utils"
)

func main() {
	utils.Init()
	r := router.Router()
	err := r.Run(":8081")
	if err != nil {
		return
	}
}
