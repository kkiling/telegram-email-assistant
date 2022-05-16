package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigApp struct {
	FileDirectory      string `yaml:"file_directory"`
	MaxTextMessageSize int    `yaml:"max_text_message_size"`
	MarkAsReadMessages bool   `yaml:"mark_as_read_messages"`
	MailCheckTimeout   int    `yaml:"mail_check_timeout"`
}

type ConfigImap struct {
	ImapServer string `yaml:"imap_server"`
	Login      string `yaml:"login"`
	Password   string `yaml:"password"`
}

type ConfigTelegram struct {
	BotToken      string `yaml:"bot_token"`
	AllowedUserId uint32 `yaml:"allowed_user_id"`
}

type Config struct {
	App      ConfigApp  `yaml:"app"`
	Imap     ConfigImap `yaml:"imap"`
	Telegram ConfigImap `yaml:"telegram"`
}

func NewConfig() (*Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("failure upload yaml file. err %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	return &config, nil
}
