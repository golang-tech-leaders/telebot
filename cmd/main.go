package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

//точка входа программы
func main() {

	token := os.Getenv("TOKEN")
	if token == "" {
		token = telToken
	}
	//https://api.telegram.org/bot<token>/METHOD_NAME
	url := telApi + token
	offset := 0

	for {
		updates, err := getUpdates(url, offset)
		if err != nil {
			log.Println("Something go wrong: ", err.Error())
		}

		sort.Slice(updates, func(i, j int) bool { return updates[i].UpdateId < updates[j].UpdateId })
		for _, update := range updates {
			if update.Message.Text == "button" {
				err := sendButton(url, update)
				if err != nil {
					log.Println("Something go wrong: ", err.Error())
				}
			} else {
				err = respond(url, update)
				if err != nil {
					log.Println("Something go wrong: ", err.Error())
				}
			}
			offset = update.UpdateId + 1
		}

		fmt.Println(updates)
	}
}

//запрос обновлений
func getUpdates(url string, offset int) ([]Update, error) {
	resp, err := http.Get(url + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

//ответ на обновления
func respond(url string, update Update) error {

	botMessage := BotMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   update.Message.Text,
	}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	resp, err := http.Post(url+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return errors.New(resp.Status)
	}

	return nil
}

func sendButton(url string, update Update) error {

	button := KeyboardButton{"не нажимайте эту кнопку"}
	keyboard := make([][]KeyboardButton, 1)
	keyboard[0] = make([]KeyboardButton, 1)
	keyboard[0][0] = button
	myKeyboard := ReplyKeyboardMarkup{
		Keyboard: keyboard,
	}
	greetingMessage := BotMessage{
		ChatId:      update.Message.Chat.ChatId,
		Text:        update.Message.Text,
		ReplyMarkup: myKeyboard,
	}

	bufGreeting, err := json.Marshal(greetingMessage)

	if err != nil {
		return err
	}

	resp, err := http.Post(url+"/sendMessage", "application/json", bytes.NewBuffer(bufGreeting))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return errors.New(resp.Status)
	}

	return nil
}
