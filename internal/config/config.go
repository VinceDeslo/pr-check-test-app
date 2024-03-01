package config

type Config struct {
	AppID string `env:"APP_ID,required"`
	WebhookSecret string `env:"WEBHOOK_SECRET,required"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH,required"`
}
