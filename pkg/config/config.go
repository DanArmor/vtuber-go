package config

import "github.com/spf13/viper"

type Config struct {
	Ip              string `mapstructure:"ip"`
	Port            string `mapstructure:"port"`
	BasePath        string `mapstructure:"base_path"`
	DriverName      string `mapstructure:"driver_name"`
	SqlUrl          string `mapstructure:"sql_url"`
	TgBotToken      string `mapstructure:"tg_bot_token"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

func LoadConfig(configPath string) (c Config, err error) {
	viper.SetConfigFile(configPath)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
