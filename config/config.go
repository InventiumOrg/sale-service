package config

import "github.com/spf13/viper"

type Config struct {
	ESEndpoint string `mapstructure:"ES_ENDPOINT"`
	ESAPIKey   string `mapstructure:"ES_API_KEY"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ClerkKey   string `mapstructure:"CLERK_KEY"`
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
	viper.Unmarshal(&config)
	return config, nil
}
