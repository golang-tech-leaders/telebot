package models

type Update struct {
	UpdateId int           `json:"update_id"`
	Message  UpdateMessage `json:"message"`
}

type UpdateMessage struct {
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Location Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"longitude"`
	Lat float64 `json:"latitude"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type GetUpdateResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId      int         `json:"chat_id"`
	Text        string      `json:"text"`
	ReplyMarkup interface{} `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type KeyboardButton struct {
	TextButton      string `json:"text"`
	RequestLocation bool   `json:"request_location"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
}

type WasteType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MapReference struct {
	Url string `json:"url"`
}

type Error struct {
	Code    int     `json:"code"`
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
}
