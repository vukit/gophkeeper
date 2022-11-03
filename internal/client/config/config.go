package config

// Config структура конфигурации приложения
type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS"`
	ServerProtocol string `env:"SERVER_PROTOCOL"`
	LogFile        string `env:"CLIENT_LOG_FILE"`
	UserInterface  string `env:"CLIENT_USER_INTERFACE"`
	DownloadFolder string `env:"CLIENT_DOWNLOAD_FOLDER"`
}
