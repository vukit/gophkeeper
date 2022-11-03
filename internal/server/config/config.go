package config

// Config структура конфигурации сервера
type Config struct {
	Address        string `env:"ADDRESS"`
	Protocol       string `env:"PROTOL"`
	TLSCertificate string `env:"TLS_CERTIFICATE"`
	TLSPrivateKey  string `env:"TLS_PRIVATE_KEY"`
	DataBaseURI    string `env:"DATABASE_URI"`
	FileStorage    string `env:"FILE_STORAGE"`
}
