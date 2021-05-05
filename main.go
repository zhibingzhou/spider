package main

import (
	"test/model"
	"test/router"
)

func main() {
	model.Delcash()
	//go proall.NewProcessor().Register().Boot()
	router.Router.Run(":8082")
}
 