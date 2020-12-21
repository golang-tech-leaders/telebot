package main

import (
	"flag"
	"fmt"
	"net/http"
	"telebot/internal/database"
	"telebot/internal/models"
	"telebot/internal/processor"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	config := prepareConfig()

	storage := database.NewTelebotLanguageStorage(config)
	storage.Migrate()
	go http.ListenAndServe("0.0.0.0:"+config.Port, nil)
	processor := processor.NewProcessor()
	processor.Start(config, storage)
}

func prepareConfig() *models.Config {
	var cfg models.Config
	configFile := getConfigFile()

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		fmt.Printf("Unable to get app configuration due to: %s\n", err.Error())
	}

	// if err := cleanenv.ReadEnv(&cfg); err != nil {
	// 	fmt.Printf("Unable to retrieve app configuration due to: %s\n", err.Error())
	// 	os.Exit(1)
	// }
	return &cfg
}

func getConfigFile() string {
	configFile := flag.String("config", "config.yml", "config file")
	return *configFile
}
