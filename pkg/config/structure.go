package config

import (
	launcherConfig "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/config"
)

const (
	DefaultLogLevel = "info"
)

type configuration struct {
	launcherConfig.StandardConfig `mapstructure:",squash"`

	Environment  Environment                  `mapstructure:"environment"`
	FuYou        FuYouSetting                 `mapstructure:"fuyou"`
	CallbackHost string                       `mapstructure:"callbackHost"`
	RPCServices  map[string]rpcServiceSetting `json:"rpcServices" yaml:"rpcServices"`
}

func (config *configuration) GetEnvironment() Environment {

	if config.Environment != "" {
		return config.Environment
	}

	return Release
}
func (config *configuration) GetCallbackHost() string {

	return config.CallbackHost
}

var (
	Config *configuration
)

func init() {
	Config = &configuration{}
}
