package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"template/pkg/api/dtos"
	token "template/pkg/helpers"
	"template/pkg/util"
)

type PUserController struct {
	Router *gin.RouterGroup
}

// 初始化路由
func (controller PUserController) Setup() {
	handle := controller.Router.Group("user")
	handle.POST("currentUser", controller.currentUser)
}

func (controller *PUserController) currentUser(c *gin.Context) {
	address, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error_code": util.ErrorCodeInternalError, "error_msg": "internal error", "data": nil})
		return
	}
	token := token.ExtractToken(c)
	c.JSON(http.StatusOK, gin.H{"error_code": util.ErrorCodeSuccess, "error_msg": "success", "data": dtos.LoginResp{
		Address: address,
		Token:   token,
	}})
}
