package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"template/pkg/api/dtos"
	"template/pkg/dao"
	"template/pkg/util"
)

type PUserController struct {
	Router *gin.RouterGroup
}

// 初始化路由
func (controller PUserController) Setup() {
	handle := controller.Router.Group("user")
	handle.POST("userInfo", controller.getUserInfo)
}

func (controller *PUserController) getUserInfo(c *gin.Context) {
	var req dtos.UserDto
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"error_code": util.ErrorCodeInvalidParams, "error_msg": "invalid params", "data": nil})
		return
	}
	if loginResp, ret := dao.SaveOrUpdateUser(req); ret {
		c.JSON(http.StatusOK, gin.H{"error_code": util.ErrorCodeSuccess, "error_msg": "success", "data": loginResp})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error_code": util.ErrorCodeInternalError, "error_msg": "internal error", "data": nil})
}
