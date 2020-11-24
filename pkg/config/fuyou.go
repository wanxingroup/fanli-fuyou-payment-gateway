package config

type FuYouSetting struct {
	OrganizationCode string `mapstructure:"organizationCode"` // 我司在富友的机构号
	Host             string `mapstructure:"host"`
	PrivatePem       string `mapstructure:"privatePemPath"`
	PublicPem        string `mapstructure:"publicPemPath"`
}

func (setting FuYouSetting) GetOrganizationCode() string {
	return setting.OrganizationCode
}

func (setting FuYouSetting) GetHost() string {
	return setting.Host
}

func (setting FuYouSetting) GetPrivatePemPath() string {
	return setting.PrivatePem
}

func (setting FuYouSetting) GetPublicPemPath() string {
	return setting.PublicPem
}
