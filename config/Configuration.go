package config

import (
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

var App *Configuration

type Configuration struct {
	Server    ServerConfiguration
	Database  DatabaseConfiguration
	Scheduler SchedulerConfiguration
	Mailer    MailerConfiguration
}

func Init() error {
	App = &Configuration{}
	viper.SetConfigName("default")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		glog.Fatalf("Failed to load default configuration file: %s", err)
		return err
	}

	viper.SetConfigName(".env")
	if err := viper.MergeInConfig(); err != nil {
		glog.Warningf("Failed to load custom configuration from .env file: %s", err)
	}

	cfg := new(Configuration)
	if err := viper.Unmarshal(cfg); err != nil {
		glog.Fatalf("Failed to deserialize config struct: %s", err)
		return err
	}
	App = cfg
	return nil
}
