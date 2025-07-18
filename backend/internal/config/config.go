package config

import "github.com/spf13/viper"

type Config struct {
	IP              string `mapstructure:"ip"`
	Port            string `mapstructure:"port"`
	BasePath        string `mapstructure:"base_path"`
	DriverName      string `mapstructure:"driver_name"`
	SQLURL          string `mapstructure:"sql_url"`
	TgBotToken      string `mapstructure:"tg_bot_token"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
	JwtSecretKey    string `mapstructure:"jwt_secret_key"`
	AdminToken      string `mapstructure:"admin_token"`
	HolodexAPIKey   string `mapstructure:"holodex_api_key"`
	IsDebug         bool   `mapstructure:"debug"`
	TimeNotifyAfter int    `mapstructure:"time_notify_after"`
	TimeStep        int    `mapstructure:"time_step"`
}

func LoadConfig(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return Config{}, err
	}

	c := Config{}
	err = viper.Unmarshal(&c)

	return c, err
}
