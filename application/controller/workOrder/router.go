package workOrder

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.GET("/query", Query)
	router.GET("/queryByContractId", QueryByContractId)
	router.GET("/queryByDrawingID", QueryByDrawingID)
	router.POST("/delete", Delete)
	router.POST("/set", Set)
	router.POST("/update", Update)
	router.POST("/queryworkorderfile", Queryworkorderfile)
}
