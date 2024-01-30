package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"os"
)

const (
	configType = "yaml"
	configName = "spec"

	// Env variable keys
	configPathEnvKey      = "SPEC_FILE_PATH"
	appEnvironment        = "APP_ENVIRONMENT"
	appIsLocalDevelopment = "APP_IS_LOCAL_DEVELOPMENT"
)

type Config struct {
	App struct {
		Name               string `mapstructure:"name" validate:"required"`
		Environment        string `mapstructure:"environment" validate:"required"`
		Timeout            int    `mapstructure:"timeout" validate:"required"`
		IsLocalDevelopment bool   `mapstructure:"is_local_development"`
	} `mapstructure:"app" validate:"required"`
}

func GetConfig() (*Config, error) {
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		return nil, fmt.Errorf("config path Env variable is not set for key, %v", configPathEnvKey)
	}

	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	err := viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error while reading config file, %v", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	// override config with env variables for deployments
	appEnv, present := os.LookupEnv(appEnvironment)

	if appEnv != "" && present {
		config.App.Environment = appEnv
	}

	_, present = os.LookupEnv(appIsLocalDevelopment)

	if present {
		if !viper.GetBool(appIsLocalDevelopment) {
			config.App.IsLocalDevelopment = false
		}
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("missing required attributes, %v", err)
	}

	return config, nil
}
