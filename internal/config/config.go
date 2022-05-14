package config

type Config struct {
	FileStorageDir string
}

func NewConfig() *Config {
	return &Config{
		// TODO: FileStorageDir
		FileStorageDir: "/home/kiling/email-data",
	}
}
