package processor

import (
	"log"
	"sort"
	"strings"
	"telebot/internal/models"
	"telebot/internal/services"
)

func ProcessUpdates(token string, offset int, userSessions map[int]models.WasteType) int {
	updates, err := services.GetUpdates(token, offset)
	if err != nil {
		return offset
	}

	sort.Slice(updates, func(i, j int) bool { return updates[i].UpdateId < updates[j].UpdateId })

	for _, update := range updates {
		messageText := strings.ToLower(update.Message.Text)
		chatId := update.Message.Chat.ChatId
		location := update.Message.Location

		if messageText == "/start" {
			processStart(token, chatId)
		} else if messageText == "" && location.Lon != 0 && location.Lat != 0 {
			if wasteType, ok := userSessions[chatId]; ok {
				if processLocation(token, chatId, location.Lat, location.Lon, wasteType.Id) {
					delete(userSessions, update.Message.Chat.ChatId)
				}
			}
		} else if messageText == "/getwastetypes" {
			processWasteTypesRequest(token, chatId)
		} else {
			processFreeText(token, chatId, messageText, userSessions)
		}
		offset = update.UpdateId + 1
	}

	return offset
}

func processStart(token string, chatId int) {
	err := services.SendTextMessage(token, chatId, "Бла-бла. Такие-то команды. Пользуйтесь на здоровье планеты.")
	if err != nil {
		services.SendTextMessage(token, chatId, err.Error())
	}
}

func processLocation(token string, chatId int, lat float64, lon float64, wasteTypeId int) bool {
	geoUrl, err := services.GetGeoUrl(wasteTypeId, lat, lon)
	if err != nil {
		services.SendTextMessage(token, chatId, err.Error())
		return false
	}

	err = services.SendTextMessage(token, chatId, geoUrl)
	if err != nil {
		services.SendTextMessage(token, chatId, err.Error())
		return false
	}

	return true
}

func processWasteTypesRequest(token string, chatId int) {
	wasteTypeList, err := services.GetWasteTypes()
	if err != nil {
		services.SendTextMessage(token, chatId, err.Error())
		return
	}

	wasteTypeNameList := make([]string, len(wasteTypeList))
	for i := range wasteTypeList {
		wasteTypeNameList[i] = wasteTypeList[i].Name
	}

	err = services.SendTextButtons(token, chatId, "Выберите отход:", wasteTypeNameList)
	if err != nil {
		log.Println("Something go wrong: ", err.Error())
	}
}

func processFreeText(token string, chatId int, text string, userSessions map[int]models.WasteType) {
	wasteType, err := services.GetWasteTypeByText(text)
	if err != nil {
		services.SendTextMessage(token, chatId, err.Error())
		return
	}

	if wasteType != (models.WasteType{}) && wasteType != nil {
		userSessions[chatId] = wasteType
		err = services.SendLocatonRequest(token, chatId, "Для получения пунктов сдачи необходимо определить геолокацию. Нажмите на кнопку:")
		if err != nil {
			services.SendTextMessage(token, chatId, err.Error())
		}
	} else {
		err = services.SendTextMessage(token, chatId, "Введенный текст не распознался как вид отхода. Можете повторить запрос либо воспользоваться командой /getwastetypes для вывода кнопок с видами отходов.")
		if err != nil {
			services.SendTextMessage(token, chatId, err.Error())
		}
	}
}
