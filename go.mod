module dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway

go 1.13

require (
	dev-gitlab.wanxingrowth.com/fanli/card v0.0.0-20200825014026-eefad6115612
	dev-gitlab.wanxingrowth.com/fanli/order/v2 v2.0.1
	dev-gitlab.wanxingrowth.com/fanli/payment v0.0.0-20200824085623-29e00ca7249d
	dev-gitlab.wanxingrowth.com/wanxin-go-micro/base v0.2.19
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/beevik/etree v1.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ozzo/ozzo-validation/v4 v4.2.2
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.12
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	golang.org/x/text v0.3.3
	google.golang.org/grpc v1.30.0
)

replace dev-gitlab.wanxingrowth.com/fanli/card => github.com/wanxingroup/fanli-card v0.0.0

replace dev-gitlab.wanxingrowth.com/wanxin-go-micro/base => github.com/wanxin-go-micro/base v0.2.27

replace dev-gitlab.wanxingrowth.com/fanli/order/v2 => github.com/wanxingroup/fanli-order v2.0.13

replace dev-gitlab.wanxingrowth.com/fanli/payment => github.com/wanxingroup/fanli-payment v0.0.0
