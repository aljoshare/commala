package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ReportJunitPath      string
	AuthorEmailEnabled   bool
	AuthorEmailWhitelist []string
	AuthorNameEnabled    bool
	AuthorNameWhitelist  []string
	BranchEnabled        bool
	BranchWhitelist      []string
	MessageEnabled       bool
	MessageWhitelist     []string
	SignOffEnabled       bool
	SignOffWhitelist     []string
}

func (c *Config) ReadConfig() {
	configPath := viper.GetString("config")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName(".commala")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.commala")
		viper.AddConfigPath(".")
	}
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
	c.AuthorEmailWhitelist = viper.GetStringSlice("validate.author.email.whitelist")
	c.AuthorNameEnabled = viper.GetBool("validate.author.name.enabled")
	c.AuthorNameWhitelist = viper.GetStringSlice("validate.author.name.whitelist")
	c.BranchEnabled = viper.GetBool("validate.branch.enabled")
	c.BranchWhitelist = viper.GetStringSlice("validate.branch.whitelist")
	c.MessageEnabled = viper.GetBool("validate.message.enabled")
	c.MessageWhitelist = viper.GetStringSlice("validate.message.whitelist")
	c.SignOffEnabled = viper.GetBool("validate.signoff.enabled")
	c.SignOffWhitelist = viper.GetStringSlice("validate.signoff.whitelist")
}
