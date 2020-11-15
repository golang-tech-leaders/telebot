package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"telebot/internal/http_error"
	"telebot/internal/models"
)

// GetWasteTypeByText Получение типа отхода из текста
func GetWasteTypeByText(text string) (*models.WasteType, error) {
	requestUrl := utils.RecyclingApiUrl + "waste/type/search/" + text
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl)
	}
	defer resp.Body.Close()

	bodyJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl)
	}

	if http_error.IsFailStatus(resp.StatusCode) {
		return nil, http_error.HttpErrorWithCustom404(resp.StatusCode,
			resp.Status+" "+requestUrl+" "+string(bodyJson),
			"Введенный текст не распознался как вид отхода. Можете повторить запрос либо воспользоваться командой /getwastetypes для вывода кнопок с видами отходов.")
	}

	var wasteType *models.WasteType
	err = json.Unmarshal(bodyJson, &wasteType)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl + " " + string(bodyJson))
	}

	return wasteType, nil
}

// GetWasteTypes Получение списка отходов
func GetWasteTypes() ([]models.WasteType, error) {
	requestUrl := utils.RecyclingApiUrl + "waste/type/list"
	resp, err := http.Get(utils.RecyclingApiUrl + "waste/type/list")
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

	var wasteTypes []models.WasteType
	err = json.Unmarshal(bodyJson, &wasteTypes)
	if err != nil {
		return nil, utils.CommonError(err.Error() + " " + requestUrl + " " + string(bodyJson))
	}

	return wasteTypes, nil
}
