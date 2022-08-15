package main

import (
	"raft-fabric-project/application/middle_ware"
	all_router "raft-fabric-project/application/router"
	"github.com/gin-gonic/gin"
)

//import "fmt"
/*
gin框架主入口

*/

func main() {
	router := gin.Default()

	router.Use(middle_ware.CorsMiddleWare)
	all_router.InitRouter(router)

	router.Run()
}
