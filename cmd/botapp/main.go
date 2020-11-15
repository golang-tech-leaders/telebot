package main

import (
	"os"
	"telebot/internal/models"
	"telebot/internal/processor"
	"telebot/internal/utils"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		token = utils.TelegramToken
	}

	offset := 0

	userSessions := make(map[int]models.WasteType)

	for {
		offset = processor.ProcessUpdates(token, offset, userSessions)
	}
}
