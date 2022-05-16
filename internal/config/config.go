package config

type Config struct {
	FileStorageDir     string
	MaxTextMessageSize int
}

func NewConfig() *Config {
	return &Config{
		// TODO: FileStorageDir
		FileStorageDir:     "/home/kiling/email-data",
		MaxTextMessageSize: 32,
	}
}
