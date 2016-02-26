package app

type AppConfig struct {
	AVKEndpoint string
	ESBEndpoint string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}