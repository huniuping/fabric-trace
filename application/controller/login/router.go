package login

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.GET("/userLoginG", UserLoginG)
	router.GET("/userRegisterG", UserRegisterG)
	router.GET("/changePasswordG", ChangePasswordG)
	router.GET("/deleteUsernameG", DeleteUsernameG)

}
