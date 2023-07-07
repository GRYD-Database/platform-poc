package configuration

import (
	"github.com/spf13/viper"
)

type Config struct {
	Address  string `mapstructure:"ADDRESS"`
	Postgres struct {
		Host       string `mapstructure:"DB_HOST"`
		Password   string `mapstructure:"DB_PASSWORD"`
		Port       string `mapstructure:"DB_PORT"`
		DBName     string `mapstructure:"DB_NAME"`
		DBUsername string `mapstructure:"DB_USERNAME"`
	} `mapstructure:"PG"`
	IPFS struct {
		RepoPath   string `mapstructure:"REPOPATH"`
		IsLocal    bool   `mapstructure:"ISLOCAL"`
		CreateRepo bool   `mapstructure:"CREATEREPO"`
	} `mapstructure:"IPFS"`
	Logger struct {
		LogLevel string `mapstructure:"LOG_LEVEL"`
		LogEnv   string `mapstructure:"LOG_ENV"`
	} `mapstructure:"LOGGER"`
	GRYDContract Contract `mapstructure:"GRYD_CONTRACT"`
	ChainConfig  Crypto   `mapstructure:"CRYPTO"`
}

type Crypto struct {
	PrivateKey string `mapstructure:"PRIVATE_KEY"`
	Endpoint   string `mapstructure:"ENDPOINT"`
}

type Contract struct {
	ABI     interface{} `mapstructure:"ABI"`
	Address string      `mapstructure:"ADDRESS"`
}

func Init() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("env.json")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = viper.Unmarshal(&config)

	return &config, nil
}
