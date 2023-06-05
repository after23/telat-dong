package util

import "github.com/spf13/viper"

type Config struct {
	Token string `mapstructure:"BOT_TOKEN"`
	MyID string `mapstructure:"MY_ID"`
	PlaygroundID string `mapstructure:"PLAYGROUND_ID"`
	ImageDumpID string `mapstructure:"IMAGE_DUMP_ID"`
	NFGuildID string `mapstructure:"NF_GUILD_ID"`
}

var conf Config

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
	err = viper.Unmarshal(&conf)
	return
}

func Conf() Config{
	return conf
}