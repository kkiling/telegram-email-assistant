package config

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type App struct {
	FileDirectory      string `yaml:"file_directory"`
	SQLiteDbFile       string `yaml:"store_db"`
	MaxTextMessageSize int    `yaml:"max_text_message_size"`
	MailCheckTimeout   int    `yaml:"mail_check_timeout"`
}

type Imap struct {
	ImapServer string `yaml:"imap_server"`
	Login      string `yaml:"login"`
	Password   string `yaml:"password"`
}

type Telegram struct {
	BotToken       string  `yaml:"bot_token"`
	AllowedUserIds []int64 `yaml:"allowed_user_id"`
}

type Config struct {
	App      App      `yaml:"app"`
	Imap     Imap     `yaml:"imap"`
	Telegram Telegram `yaml:"telegram"`
}

func NewConfig(configFile string) (*Config, error) {
	logrus.Infof("Read config file: %s", configFile)
	var config Config
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failure upload yaml file. err %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return &config, nil
}
