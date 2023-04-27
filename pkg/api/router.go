package api

import (
	"github.com/gin-gonic/gin"
	"template/pkg/middleware"
)

// Setup sets up all controllers.
func Setup(router *gin.Engine) {
	public := router.Group("api")
	user := UserController{Router: public}
	user.Setup()

	protected := router.Group("api")
	protected.Use(middleware.JwtAuthMiddleware())

	puser := PUserController{Router: protected}
	puser.Setup()

}
