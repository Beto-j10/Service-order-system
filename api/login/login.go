package login

import "github.com/gin-gonic/gin"

// Login authenticates the user
func Login(router *gin.Engine) {
	initMiddleware()
	router.POST("/login", AuthMiddleware.LoginHandler)
}
