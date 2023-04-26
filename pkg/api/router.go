package api

import "github.com/gin-gonic/gin"

// Setup sets up all controllers.
func Setup(router *gin.Engine) {
	api := router.Group("api")

	user := UserController{Router: api}
	user.Setup()
}
