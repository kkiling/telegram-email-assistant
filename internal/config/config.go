package config

type Config struct {
	FileStorageDir     string
	MaxTextMessageSize int
	ImapServer         string
	Login              string
	Password           string
}

func NewConfig() *Config {
	return &Config{
		// TODO: FileStorageDir
	}
}
