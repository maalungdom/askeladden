package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord struct {
		Token          string `yaml:"token"`
		Prefix         string `yaml:"prefix"`
		LogChannelID   string `yaml:"logChannelID"`
		DefaultChannelID string `yaml:"defaultChannelID"`
	} `yaml:"discord"`

	Approval struct {
		QueueChannelID string `yaml:"queueChannelID"`
		OpplysarRoleID string `yaml:"opplysarRoleID"`
	} `yaml:"approval"`

	Starboard struct {
		ChannelID string `yaml:"channelID"`
		Threshold int    `yaml:"threshold"`
	} `yaml:"starboard"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Scheduler struct {
		CronString string `yaml:"cron_string"`
	} `yaml:"scheduler"`

	// WarningEmbedConfig enables or disables with an optional custom text
	ShowAISlopWarning bool   `yaml:"showAISlopWarning"`
	AISlopWarningText string `yaml:"aiSlopWarningText"`
}


// FUNKSJON. Lastar inn konfigurasjonen og gir ein fylt Config-struct
//--------------------------------------------------------------------------------
func Load() (*Config, error) {
	var cfg Config

//1. Opnar config.yaml og les ho.
//--------------------------------------------------------------------------------
	configFile, err := os.Open("config.yaml")
	if err != nil {	
		return nil, err
	}
	defer configFile.Close()
	
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
//2. Opnar secrets.yaml og les ho.
//--------------------------------------------------------------------------------
	secretsFile, err := os.Open("secrets.yaml")
	if err != nil {
		return nil, err
	}
	defer secretsFile.Close()

	secretsDecoder := yaml.NewDecoder(secretsFile)
	if err := secretsDecoder.Decode(&cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}
