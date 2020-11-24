package state

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
	"github.com/gin-gonic/gin"
)

// @ID Ping
// @Summary Test service state
// @Description Just test service is running
// @Tags ping
// @Success 200 {string} string	"PONG"
// @Router /ping [get]
func (_ Controller) Ping(ctx *gin.Context) {

	response.Response(ctx, "PONG")
}
