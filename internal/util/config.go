package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	GinMode           string `mapstructure:"GIN_MODE"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	MigrationSrc      string `mapstructure:"MIGRATION_SRC"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	APIBasePath       string `mapstructure:"API_BASE_PATH"`
	OutputPath        string `mapstructure:"OUTPUT_PATH"`
	WebDistPath       string `mapstructure:"WEB_DIST_PATH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
