package design

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.GET("/query", Query)
	router.GET("/queryByContractId", QueryByContractId)
	router.GET("/queryDrawingDetailByID", QueryDrawingDetailByID)
	router.GET("/querydrawingfile", Querydrawingfile)
	router.POST("/delete", Delete)
	router.POST("/set", Set)
	router.POST("/update", Update)
}
