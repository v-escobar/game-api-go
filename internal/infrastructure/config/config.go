package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Db struct {
		Host     string `mapstructure:"host"`
		Port     uint32 `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}

	App struct {
		Port uint `mapstructure:"port"`
	}
}

func (app *AppConfig) Dsn() string {
	return fmt.Sprint(
		"host=", app.Db.Host,
		" user=", app.Db.User,
		" password=", app.Db.Password,
		" port=", app.Db.Port,
		" database=", app.Db.Database,
	)
}

func InitConfig() AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("conf")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config AppConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling configuration, %s", err)
	}

	return config
}
