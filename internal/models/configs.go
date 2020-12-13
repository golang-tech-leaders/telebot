package models

type Config struct {
	TelegramToken   string `yaml:"telegram_token" env:"TELEGRAM_TOKEN"`
	TelegramApiUrl  string `yaml:"telegram_api_url" env:"TELEGRAM_API_URL" env-default:"https://api.telegram.org/bot"`
	RecyclingApiUrl string `yaml:"recycling_api_url" env:"RECYCLING_API_URL" env-default:"https://virtserver.swaggerhub.com/Trepka/mock/1.0.0/"`
	GeobaseApiUrl   string `yaml:"geobase_api_url" env:"GEOBASE_API_URL" env-default:"https://virtserver.swaggerhub.com/Trepka/Geobase2/1.0.0/"`
}
