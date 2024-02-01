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
	configPathEnvKey              = "SPEC_FILE_PATH"
	appEnvironment                = "APP_ENVIRONMENT"
	appIsLocalDevelopment         = "APP_IS_LOCAL_DEVELOPMENT"
	dynamodbTableNameMembership   = "DYNAMODB_TABLE_NAME_MEMBERSHIP"
	webhookKeyMembershipStarted   = "WEB_HOOK_KEY_MEMBERSHIP_STARTED"
	webhookKeyMembershipUpdated   = "WEB_HOOK_KEY_MEMBERSHIP_UPDATED"
	webhookKeyMembershipCancelled = "WEB_HOOK_KEY_MEMBERSHIP_CANCELLED"
)

type Config struct {
	App struct {
		Name               string `mapstructure:"name" validate:"required"`
		Environment        string `mapstructure:"environment" validate:"required"`
		Timeout            int    `mapstructure:"timeout" validate:"required"`
		IsLocalDevelopment bool   `mapstructure:"is_local_development"`
		WebHookKeys        struct {
			MembershipStarted   string `mapstructure:"membership_started" validate:"required"`
			MembershipUpdated   string `mapstructure:"membership_updated" validate:"required"`
			MembershipCancelled string `mapstructure:"membership_cancelled" validate:"required"`
		} `mapstructure:"web_hook_keys" validate:"required"`
	} `mapstructure:"app" validate:"required"`
	Aws struct {
		Dynamodb struct {
			Membership struct {
				TableName string `mapstructure:"table_name" validate:"required"`
			} `mapstructure:"membership" validate:"required"`
		} `mapstructure:"dynamodb" validate:"required"`
	} `mapstructure:"aws" validate:"required"`
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
	viper.AutomaticEnv()
	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	err = overrideConfigWithEnvVariables(config)

	if err != nil {
		return nil, fmt.Errorf("overriedConfigWithEnvVariables error: %v", err)
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("missing required attributes, %v", err)
	}

	return config, nil
}

func overrideConfigWithEnvVariables(config *Config) error {
	// override config with env variables for deployments
	appEnv, present := os.LookupEnv(appEnvironment)

	if appEnv != "" && present {
		config.App.Environment = appEnv
	} else {
		return fmt.Errorf("critical env variable is missing")
	}

	_, present = os.LookupEnv(appIsLocalDevelopment)

	if present {
		if !viper.GetBool(appIsLocalDevelopment) {
			config.App.IsLocalDevelopment = false
		}
	} else {
		return fmt.Errorf("critical env variable is missing")
	}

	membershipStartedSecretKey, present := os.LookupEnv(webhookKeyMembershipStarted)

	if membershipStartedSecretKey != "" && present {
		config.App.WebHookKeys.MembershipStarted = membershipStartedSecretKey
	} else {
		return fmt.Errorf("critical env variable is missing")
	}

	membershipUpdatedSecretKey, present := os.LookupEnv(webhookKeyMembershipUpdated)

	if membershipUpdatedSecretKey != "" && present {
		config.App.WebHookKeys.MembershipUpdated = membershipUpdatedSecretKey
	} else {
		return fmt.Errorf("critical env variable is missing")
	}

	membershipCancelledSecretKey, present := os.LookupEnv(webhookKeyMembershipCancelled)

	if membershipCancelledSecretKey != "" && present {
		config.App.WebHookKeys.MembershipCancelled = membershipCancelledSecretKey
	} else {
		return fmt.Errorf("critical env variable is missing")
	}

	tableNameMembership := viper.GetString(dynamodbTableNameMembership)

	if tableNameMembership != "" {
		config.Aws.Dynamodb.Membership.TableName = tableNameMembership
	}

	return nil
}
