package models

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	SES      SESConfig
	S3       S3Config
	Telegram TelegramConfig
	RabbitMQ RabbitMQConfig
	Cors     CorsConfig
}

type CorsConfig struct {
	AllowedOrigins []string
}

type ServerConfig struct {
	Host         string
	Port         string
	PortFrontend string
	KeyPassword  string
}

type DatabaseConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
}

type CacheConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type SESConfig struct {
	Region string
	Sender string // verified sender address in SES
}

type S3Config struct {
	Region string
	Bucket string
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

type RabbitMQConfig struct {
	Username string
	Password string
	URL      string
}
