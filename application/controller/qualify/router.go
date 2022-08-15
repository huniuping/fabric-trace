package quality

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.GET("/query", Query)
	router.GET("/queryQulityDetailByID", QueryQulityDetailByID)
	router.POST("/delete", Delete)
	router.POST("/set", Set)
	router.POST("/update", Update)
}
