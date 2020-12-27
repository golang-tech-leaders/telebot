package processor

import (
	"fmt"
	"log"
	"strings"
	"telebot/internal/database"
	"telebot/internal/http_error"
	"telebot/internal/models"
	"telebot/internal/services"
)

type Processor struct {
	closeChan chan int
}

func NewProcessor() *Processor {
	return &Processor{
		closeChan: make(chan int),
	}
}

func (p *Processor) Start(config *models.Config, storage *database.TelebotLanguageStorage) {
	offset := 0
	updatesChan := make(chan []models.Update)
	userSessionsWaste := make(map[int]models.WasteType)
	userSessionsLang := make(map[int]map[int]string)

	for {
		go func() {
			updatesList, err := services.GetUpdates(config.TelegramToken, config.TelegramApiUrl, offset)
			if err != nil {
				return
			}
			updatesChan <- updatesList
		}()

		select {
		case <-p.closeChan:
			return
		case updatesList := <-updatesChan:
			for _, update := range updatesList {
				messageText := strings.ToLower(update.Message.Text)
				chatId := update.Message.Chat.ChatId
				location := update.Message.Location

				if userSessionsLang[chatId] == nil {
					messageMap, err := storage.GetLangMessage(database.RU)
					if err != nil {
						fmt.Println(err)
					}
					userSessionsLang[chatId] = *messageMap
				}

				switch messageText {
				case "/start":
					processStart(config, chatId, userSessionsLang[chatId][1])
				case "/getwastetypes":
					processWasteTypesRequest(config, chatId, userSessionsLang[chatId][3])
				case "/en":
					messageMap, err := storage.GetLangMessage(database.EN)
					if err != nil {
						fmt.Println(err)
					}
					userSessionsLang[chatId] = *messageMap
				case "/ru":
					messageMap, err := storage.GetLangMessage(database.RU)
					if err != nil {
						fmt.Println(err)
					}
					userSessionsLang[chatId] = *messageMap
				default:
					if messageText == "" && location.Lon != 0 && location.Lat != 0 {
						if wasteType, ok := userSessionsWaste[chatId]; ok {
							if processLocation(config, chatId, location.Lat, location.Lon, wasteType.Name, userSessionsLang[chatId][2]) {
								delete(userSessionsWaste, update.Message.Chat.ChatId)
							}
						}
					} else {
						processFreeText(config, chatId, messageText, userSessionsWaste, userSessionsLang[chatId][5], userSessionsLang[chatId][4], userSessionsLang[chatId][6])
					}
				}
				offset = update.UpdateId + 1
			}
		}
	}
}

func (p *Processor) Stop() {
	close(p.closeChan)
}

func processStart(config *models.Config, chatId int, helloMsg string) {
	err := services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, helloMsg)
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
	}
}

func processLocation(config *models.Config, chatId int, lat float64, lon float64, wasteTypeId string, pointNotFound string) bool {
	geoUrl, err := services.GetGeoUrl(config.GeobaseApiUrl, wasteTypeId, lat, lon)
	if err != nil {
		if err == http_error.ErrNotFound {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, pointNotFound)
		} else {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
		return false
	}

	if geoUrl != nil {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, *geoUrl)
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
			return false
		}
	} else {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, pointNotFound)
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	}

	return true
}

func processWasteTypesRequest(config *models.Config, chatId int, chooseWasteMsg string) {
	wasteTypeList, err := services.GetWasteTypes(config.RecyclingApiUrl)
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		return
	}

	wasteTypeNameList := make([]string, len(wasteTypeList))
	for i := range wasteTypeList {
		wasteTypeNameList[i] = wasteTypeList[i].Name
	}

	err = services.SendTextButtons(config.TelegramToken, config.TelegramApiUrl, chatId, chooseWasteMsg, wasteTypeNameList)
	if err != nil {
		log.Println("Something go wrong: ", err.Error())
	}
}

func processFreeText(config *models.Config, chatId int, text string, userSessions map[int]models.WasteType, wasteNotFound string, reqLocMsg string, locBtnMsg string) {
	wasteType, err := services.GetWasteTypeByText(config.RecyclingApiUrl, text)
	if err != nil {
		if err == http_error.ErrNotFound {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, wasteNotFound)
		} else {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
		return
	}

	if wasteType != nil && *wasteType != (models.WasteType{}) {
		userSessions[chatId] = *wasteType
		err = services.SendLocatonRequest(config.TelegramToken, config.TelegramApiUrl, chatId, reqLocMsg, locBtnMsg)
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	} else {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, wasteNotFound)
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	}
}
