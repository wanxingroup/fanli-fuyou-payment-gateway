package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/client"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/cronjob"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
	cronLog "dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log/cron"
)

var c *cron.Cron

func Start() {

	app := launcher.NewApplication(
		launcher.SetApplicationDescription(
			&launcher.ApplicationDescription{
				ShortDescription: "users service",
				LongDescription:  "support mini program user data management function.",
			},
		),
		launcher.SetApplicationLogger(log.GetLogger()),
		launcher.SetApplicationEvents(
			launcher.NewApplicationEvents(
				launcher.SetOnInitEvent(func(app *launcher.Application) {

					unmarshalConfiguration()

					registerWebRouter(app)
					registerRPCRouter(app)

					idcreator.InitCreator(app.GetServiceId())

					client.InitCardService()
					client.InitOrderService()

					registerCronJobInit()
				}),
				launcher.SetOnStartEvent(func(app *launcher.Application) {

					autoMigration()
					c.Start()
				}),
				launcher.SetOnCloseEvent(func(app *launcher.Application) {

					ctx := c.Stop()
					ctx.Done()
				}),
			),
		),
	)

	app.Launch()
}

func registerCronJobInit() {

	c = cron.New(
		cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
		cron.WithLogger(cronLog.NewLogger(log.GetLogger())),
	)

	type Cron struct {
		Spec string
		Run  func()
	}

	jobs := []Cron{
		{Spec: "* * * * *", Run: cronjob.CheckTransactionsPayStatus},
	}

	for _, job := range jobs {
		_, err := c.AddFunc(job.Spec, job.Run)
		if err != nil {
			log.GetLogger().WithError(err).Error("register cron job error")
		}
	}
}

func unmarshalConfiguration() {
	err := viper.Unmarshal(config.Config)
	if err != nil {

		log.GetLogger().WithError(err).Error("unmarshal config error")
	}
}
