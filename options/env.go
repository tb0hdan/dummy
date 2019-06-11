package options

import (
	"github.com/spf13/viper"
)

func ReadEnv() *Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("PROMETHEUS_PORT", 9090)

	return &Config{
		LogLevel:        viper.GetString("LOG_LEVEL"),
		Port:            viper.GetInt("HTTP_PORT"),
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		PrometheusPort:  viper.GetInt("PROMETHEUS_PORT"),
	}
}
