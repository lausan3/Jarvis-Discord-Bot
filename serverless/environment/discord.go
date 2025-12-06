package environment

type DiscordConfig struct {
	AppID     string `mapstructure:"DISCORD_APP_ID"`
	Token     string `mapstructure:"DISCORD_TOKEN"`
	PublicKey string `mapstructure:"DISCORD_PUBLIC_KEY"`

	APIUrl string `mapstructure:"DISCORD_API_URL"`
}
