package configs

import (
	"log"
	"os"

	"go_template/configs/common/constants"
	"go_template/internal/models"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (config models.Config, err error) {
	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	viper.AddConfigPath(path)
	env := os.Getenv("ENV")
	if env == constants.DevEnvironment {
		viper.SetConfigName("config.dev")
	} else {
		viper.SetConfigName("config.prod")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Unmarshal the configuration into the config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
