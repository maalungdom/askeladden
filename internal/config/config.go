package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord struct {
		Token            string `yaml:"token"`
		Prefix           string `yaml:"prefix"`
		LogChannelID     string `yaml:"logChannelID"`
		DefaultChannelID string `yaml:"defaultChannelID"`
	} `yaml:"discord"`

	Approval struct {
		QueueChannelID string `yaml:"queueChannelID"`
		OpplysarRoleID string `yaml:"opplysarRoleID"`
	} `yaml:"approval"`

	BannedWords struct {
		ApprovalChannelID string `yaml:"approvalChannelID"`
		RettskrivarRoleID string `yaml:"rettskrivarRoleID"`
	} `yaml:"bannedwords"`

	Grammar struct {
		ChannelID string `yaml:"channelID"`
	} `yaml:"grammar"`

	Starboard struct {
		ChannelID string `yaml:"channelID"`
		Threshold int    `yaml:"threshold"`
		Emoji     string `yaml:"emoji"`
	} `yaml:"starboard"`

	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Scheduler struct {
		CronString      string `yaml:"cron_string"`
		Timezone        string `yaml:"timezone"`
		MorningTime     string `yaml:"morning_time"`
		EveningTime     string `yaml:"evening_time"`
		InactivityHours int    `yaml:"inactivity_hours"`
		Enabled         bool   `yaml:"enabled"`
	} `yaml:"scheduler"`

	// Reaction emojis
	Reactions struct {
		Question string `yaml:"question"`
	} `yaml:"reactions"`

	// Beta environment settings
	Environment       string `yaml:"environment"`
	TableSuffix       string `yaml:"table_suffix"`
	ShowAISlopWarning bool   `yaml:"showAISlopWarning"`
	AISlopWarningText string `yaml:"aiSlopWarningText"`
}

// FUNKSJON. Lastar inn konfigurasjonen og gir ein fylt Config-struct
// --------------------------------------------------------------------------------
func Load() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config/config.yaml"
	}

	secretsFile := os.Getenv("SECRETS_FILE")
	if secretsFile == "" {
		secretsFile = "config/secrets.yaml"
	}

	return LoadWithFiles(configFile, secretsFile)
}

// LoadWithFiles loads config with custom filenames
func LoadWithFiles(configFile, secretsFile string) (*Config, error) {
	var cfg Config

	// 1. Open and read config file
	cfgFile, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	decoder := yaml.NewDecoder(cfgFile)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	// 2. Open and read secrets file into a temporary struct
	secrFile, err := os.Open(secretsFile)
	if err != nil {
		return nil, err
	}
	defer secrFile.Close()

	var secrets struct {
		Discord struct {
			Token string `yaml:"token"`
		} `yaml:"discord"`
		Database struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"database"`
	}

	secretsDecoder := yaml.NewDecoder(secrFile)
	if err := secretsDecoder.Decode(&secrets); err != nil {
		return nil, err
	}

	// 3. Merge secrets into the main config
	cfg.Discord.Token = secrets.Discord.Token
	cfg.Database.User = secrets.Database.User
	cfg.Database.Password = secrets.Database.Password

	return &cfg, nil
}
