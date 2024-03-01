package config

type Config struct {
	AppID int64 `env:"APP_ID,required"`
	InstallationID int64 `env:"INSTALLATION_ID,required"`
	WebhookSecret string `env:"WEBHOOK_SECRET,required"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH,required"`
}
