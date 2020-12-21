package models

type Config struct {
	TelegramToken   string `yaml:"telegram_token" env:"TELEGRAM_TOKEN"`
	TelegramApiUrl  string `yaml:"telegram_api_url" env:"TELEGRAM_API_URL" env-default:"https://api.telegram.org/bot"`
	RecyclingApiUrl string `yaml:"recycling_api_url" env:"RECYCLING_API_URL" env-default:"https://virtserver.swaggerhub.com/Trepka/mock/1.0.0/"`
	GeobaseApiUrl   string `yaml:"geobase_api_url" env:"GEOBASE_API_URL" env-default:"https://virtserver.swaggerhub.com/Trepka/Geobase2/1.0.0/"`
	DbURL           string `yaml:"db_address" env:"DATABASE_URL" env-default:""`
	DbPort          int    `yaml:"db_port" env:"DBPORT" env-default:"5432"`
	DbHost          string `yaml:"db_host" env:"DBHOST" env-default:"localhost"`
	DbName          string `yaml:"db_name" env:"DBNAME" env-default:"postgres"`
	DbUser          string `yaml:"db_user" env:"DBUSER" env-default:"postgres"`
	DbPassword      string `yaml:"db_password" env:"DBPASSWORD"`
	Port            string `yaml:"port" env:"PORT" env-default:"8090"`
}
