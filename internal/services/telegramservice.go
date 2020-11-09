package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"telebot/internal/models"
	"telebot/internal/utils"
)

// Запрос обновлений
func GetUpdates(token string, offset int) ([]models.Update, error) {
	requestUrl := utils.TelegramApiUrl + token + "/getUpdates" + "?offset=" + strconv.Itoa(offset)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, utils.CommonError(err.Error() + " " + requestUrl)
	}
	defer resp.Body.Close()

	bodyJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.CommonError(err.Error() + " " + requestUrl)
	}

	if utils.IsFailStatus(resp.StatusCode) {
		return nil, utils.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+string(bodyJson))
	}

	var getUpdateResponse models.GetUpdateResponse
	err = json.Unmarshal(bodyJson, &getUpdateResponse)
	if err != nil {
		return nil, utils.CommonError(err.Error() + " " + requestUrl)
	}

	return getUpdateResponse.Result, nil
}

func SendTextMessage(token string, chatId int, text string) error {
	requestUrl := utils.TelegramApiUrl + token + "/sendMessage"

	message := models.BotMessage{
		ChatId:      chatId,
		Text:        text,
		ReplyMarkup: models.ReplyKeyboardRemove{RemoveKeyboard: true},
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}
	if utils.IsFailStatus(resp.StatusCode) {
		return utils.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}

func SendTextButtons(token string, chatId int, text string, textList []string) error {
	requestUrl := utils.TelegramApiUrl + token + "/sendMessage"

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
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}
	if utils.IsFailStatus(resp.StatusCode) {
		return utils.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}

func SendLocatonRequest(token string, chatId int, text string) error {
	requestUrl := utils.TelegramApiUrl + token + "/sendMessage"

	buttons := make([][]models.KeyboardButton, 1)
	keyboard := models.ReplyKeyboardMarkup{
		Keyboard: buttons,
	}
	button := models.KeyboardButton{
		TextButton:      "Отправить геолокацию",
		RequestLocation: true,
	}
	buttons[0] = make([]models.KeyboardButton, 1)
	buttons[0][0] = button

	message := models.BotMessage{
		ChatId:      chatId,
		Text:        text,
		ReplyMarkup: keyboard,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}

	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return utils.CommonError(err.Error() + requestUrl + " " + fmt.Sprintf("%v", message))
	}
	if utils.IsFailStatus(resp.StatusCode) {
		return utils.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+fmt.Sprintf("%v", message)+" "+string(messageJson))
	}

	return nil
}
