package util

import "github.com/spf13/viper"

type Config struct {
	Token string `mapstructure:"BOT_TOKEN"`
	PlaygroundID string `mapstructure:"PLAYGROUND_ID"`
	ImageDumpID string `mapstructure:"IMAGE_DUMP_ID"`
	NFGuildID string `mapstructure:"NF_GUILD_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}