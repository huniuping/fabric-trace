package contract

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.GET("/query", Query)
	router.GET("/queryByProductName", QueryByProductName)
	router.GET("/querycontractfile", Querycontractfile)
	router.GET("/queryContractDetailByID", QueryContractDetailByID)
	router.POST("/delete", Delete)
	router.POST("/set", Set)
	router.POST("/update", Update)
}
