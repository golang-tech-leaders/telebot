package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"telebot/internal/http_error"
	"telebot/internal/models"
)

// GetUpdates Запрос обновлений
func GetUpdates(token string, url string, offset int) ([]models.Update, error) {
	requestUrl := fmt.Sprintf("%s%s/getUpdates?offset=%d", url, token, offset)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl)
	}
	defer resp.Body.Close()

	if http_error.IsFailStatus(resp.StatusCode) {
		return nil, http_error.HttpError(resp.StatusCode, fmt.Sprintf("%s %s %v", resp.Status, requestUrl, resp))
	}

	var getUpdateResponse models.GetUpdateResponse
	err = json.NewDecoder(resp.Body).Decode(&getUpdateResponse)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl)
	}

	return getUpdateResponse.Result, nil
}

func SendTextMessage(token string, url string, chatId int, text string) error {
	requestUrl := url + token + "/sendMessage"

	message := models.BotMessage{
		ChatId:      chatId,
		Text:        text,
		ReplyMarkup: models.ReplyKeyboardRemove{RemoveKeyboard: true},
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return http_error.CommonError(fmt.Sprintf("%v: %s %v", err, requestUrl, message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return http_error.CommonError(fmt.Sprintf("%v: %s %v", err, requestUrl, message))
	}
	if http_error.IsFailStatus(resp.StatusCode) {
		return http_error.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}

func SendTextButtons(token string, url string, chatId int, text string, textList []string) error {
	requestUrl := url + token + "/sendMessage"

	buttons := make([][]models.KeyboardButton, len(textList))
	keyboard := models.ReplyKeyboardMarkup{
		Keyboard: buttons,
	}
	for i := range textList {
		button := models.KeyboardButton{
			TextButton:      textList[i],
			RequestLocation: false,
		}
		buttons[i] = make([]models.KeyboardButton, 1)
		buttons[i][0] = button
	}

	message := models.BotMessage{
		ChatId:      chatId,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return http_error.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return http_error.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}
	if http_error.IsFailStatus(resp.StatusCode) {
		return http_error.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}

func SendLocatonRequest(token string, url string, chatId int, text string, btnMessage string) error {
	requestUrl := url + token + "/sendMessage"

	keyboard := models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{
					TextButton:      btnMessage,
					RequestLocation: true,
				},
			},
		},
	}

	message := models.BotMessage{
		ChatId:      chatId,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return http_error.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return http_error.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}
	if http_error.IsFailStatus(resp.StatusCode) {
		return http_error.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}
