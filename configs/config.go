package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	LimitRequestsDefaultByIP int64  `mapstructure:"LIMIT_REQUESTS_DEFAULT_BY_IP"`
	RequestLimitInSec        int64  `mapstructure:"REQUEST_LIMIT_IN_SEC"`
	BlockDuration            int    `mapstructure:"BLOCK_DURATION"`
	SecretKey                string `mapstructure:"SECRET_KEY"`
	ExpirationToken          int    `mapstructure:"EXPIRATION_TOKEN"`
	LimitRequestsByToken     int64  `mapstructure:"LIMIT_REQUESTS_BY_TOKEN"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		viper.SetConfigFile("../.env")
		err = viper.ReadInConfig()
		if err != nil {
			return config, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func GetLimitRequestsDefaultByIP() int64 {
	config, _ := LoadConfig()
	return config.LimitRequestsDefaultByIP
}

func GetRequestLimitInSec() int64 {
	config, _ := LoadConfig()
	return config.RequestLimitInSec
}

func GetBlockDuration() int {
	config, _ := LoadConfig()
	return config.BlockDuration
}

func GetSecretKey() string {
	config, _ := LoadConfig()
	return config.SecretKey
}

func GetExpirationToken() int {
	config, _ := LoadConfig()
	return config.ExpirationToken
}

func GetLimitRequestsByToken() int64 {
	config, _ := LoadConfig()
	return config.LimitRequestsByToken
}
