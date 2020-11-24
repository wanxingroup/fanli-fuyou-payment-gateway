package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/restful/callback"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/restful/state"
	_ "dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/restful/swaggerdocs"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func registerWebRouter(app *launcher.Application) {
	ginService := app.GetWebService()
	if ginService == nil {

		log.GetLogger().WithField("stage", "onInit").Error("get gin service is nil")
		return
	}

	ginEngine := ginService.GetEngine()
	if ginEngine == nil {
		log.GetLogger().Error("get gin engine is nil, register web router failed")
		return
	}

	registerStateWebRouter(ginEngine)
	registerSwaggerRouter(ginEngine)
	registerCallbackRouter(ginEngine)
}

func registerCallbackRouter(ginEngine *gin.Engine) {

	controller := &callback.Controller{}
	ginEngine.POST("/api/paymentgateway/fuyou/callback", controller.Notify)
}

func registerSwaggerRouter(ginEngine *gin.Engine) {

	ginEngine.GET("/api/paymentgateway/fuyou/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func registerStateWebRouter(ginEngine *gin.Engine) {

	stateController := &state.Controller{}

	ginEngine.GET("/api/paymentgateway/fuyou/state/ping", stateController.Ping)
}
