package processor

import (
	"log"
	"sort"
	"strings"
	"telebot/internal/models"
	"telebot/internal/services"
)

type Processor struct {
	closeChan chan int
}

func NewProcessor() *Processor{
	return &Processor{
		closeChan: make(chan int),
	}	
}

type message struct {
	
}

func (p *Processor) Start(config *models.Config, userSessions map[int]models.WasteType) {
	offset := 0
	updates := make(chan []models.Update)
	
	for {
		go func() {
			u, err := services.GetUpdates(config.TelegramToken, config.TelegramApiUrl, offset)
			if err != nil {
				return
			}
			updates <- u
			
		}()

		select {
		case <-p.closeChan:
			return
		case u := <-updates:
			for upd := range u {
				messageText := strings.ToLower(update.Message.Text)
				chatId := update.Message.Chat.ChatId
				location := update.Message.Location

				switch messageText {
				case "/start":
					processStart(config, chatId)
				case "/getwastetypes":
					processWasteTypesRequest(config, chatId)	
				default:
					if messageText == "" && location.Lon != 0 && location.Lat != 0 {
						if processLocation(config, chatId, location.Lat, location.Lon, wasteType.Id) {
							delete(userSessions, update.Message.Chat.ChatId)
						}
					} else {
						processFreeText(config, chatId, messageText, userSessions)	
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

func ProcessUpdates(config *models.Config, offset int, userSessions map[int]models.WasteType) int {
	updates, err := services.GetUpdates(config.TelegramToken, config.TelegramApiUrl, offset)
	if err != nil {
		return offset
	}

	sort.Slice(updates, func(i, j int) bool { return updates[i].UpdateId < updates[j].UpdateId })

	for _, update := range updates {
		messageText := strings.ToLower(update.Message.Text)
		chatId := update.Message.Chat.ChatId
		location := update.Message.Location

		if messageText == "/start" {
			processStart(config, chatId)
		} else if messageText == "" && location.Lon != 0 && location.Lat != 0 {
			if wasteType, ok := userSessions[chatId]; ok {
				if processLocation(config, chatId, location.Lat, location.Lon, wasteType.Id) {
					delete(userSessions, update.Message.Chat.ChatId)
				}
			}
		} else if messageText == "/getwastetypes" {
			processWasteTypesRequest(config, chatId)
		} else {
			processFreeText(config, chatId, messageText, userSessions)
		}
		offset = update.UpdateId + 1
	}

	return offset
}

func processStart(config *models.Config, chatId int) {
	err := services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, "Бла-бла. Такие-то команды. Пользуйтесь на здоровье планеты.")
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
	}
}

func processLocation(config *models.Config, chatId int, lat float64, lon float64, wasteTypeId int) bool {
	geoUrl, err := services.GetGeoUrl(config.GeobaseApiUrl, wasteTypeId, lat, lon)
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		return false
	}

	if geoUrl != nil {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, *geoUrl)
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
			return false
		}
	} else {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, "Ни одного пункта сдачи не найдено.")
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	}

	return true
}

func processWasteTypesRequest(config *models.Config, chatId int) {
	wasteTypeList, err := services.GetWasteTypes(config.RecyclingApiUrl)
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		return
	}

	wasteTypeNameList := make([]string, len(wasteTypeList))
	for i := range wasteTypeList {
		wasteTypeNameList[i] = wasteTypeList[i].Name
	}

	err = services.SendTextButtons(config.TelegramToken, config.TelegramApiUrl, chatId, "Выберите отход:", wasteTypeNameList)
	if err != nil {
		log.Println("Something go wrong: ", err.Error())
	}
}

func processFreeText(config *models.Config, chatId int, text string, userSessions map[int]models.WasteType) {
	wasteType, err := services.GetWasteTypeByText(config.RecyclingApiUrl, text)
	if err != nil {
		services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		return
	}

	if wasteType != nil && *wasteType != (models.WasteType{}) {
		userSessions[chatId] = *wasteType
		err = services.SendLocatonRequest(config.TelegramToken, config.TelegramApiUrl, chatId, "Для получения пунктов сдачи необходимо определить геолокацию. Нажмите на кнопку:")
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	} else {
		err = services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, "Введенный текст не распознался как вид отхода. Можете повторить запрос либо воспользоваться командой /getwastetypes для вывода кнопок с видами отходов.")
		if err != nil {
			services.SendTextMessage(config.TelegramToken, config.TelegramApiUrl, chatId, err.Error())
		}
	}
}
