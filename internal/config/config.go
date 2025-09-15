package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ReportJunitPath    string
	AuthorEmailEnabled bool
	AuthorNameEnabled  bool
	BranchEnabled      bool
	MessageEnabled     bool
	SignOffEnabled     bool
}

func (c *Config) ReadConfig() {
	viper.SetConfigName(".commala")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.commala")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			dirname, err := os.UserHomeDir()
			if err != nil {
				log.Errorf("Can't get name of home directory")
			}
			log.Warnf("Configuration file is missing: ./config.yaml or %s/.commala/.commala.yaml", dirname)
		} else {
			log.Errorf("Can't read configuration file: %s", err)
		}
	}
	log.Debug("All configuration settings:")
	for _, key := range viper.AllKeys() {
		value := viper.Get(key)
		log.Debugf("%s: %v\n", key, value)
	}
	c.ReportJunitPath = viper.GetString("report.junit.path")
	c.AuthorEmailEnabled = viper.GetBool("validate.author.email.enabled")
	c.AuthorNameEnabled = viper.GetBool("validate.author.name.enabled")
	c.BranchEnabled = viper.GetBool("validate.branch.enabled")
	c.MessageEnabled = viper.GetBool("validate.message.enabled")
	c.SignOffEnabled = viper.GetBool("validate.signoff.enabled")
}
