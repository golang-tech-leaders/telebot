package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"telebot/internal/models"
	"telebot/internal/utils"
)

// Получение ссылки из geobase
func GetGeoUrl(wasteTypeId int, lat float64, lon float64) (string, error) {
	requestUrl := fmt.Sprintf(utils.GeobaseApiUrl+"waste/type/"+strconv.Itoa(wasteTypeId)+"/location?latitude=%f&longitude=%f", lat, lon)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", utils.CommonError(err.Error() + " " + requestUrl)
	}
	defer resp.Body.Close()

	bodyJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", utils.CommonError(err.Error() + " " + requestUrl)
	}

	if utils.IsFailStatus(resp.StatusCode) {
		return "", utils.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+string(bodyJson))
	}

	var mapReference models.MapReference
	err = json.Unmarshal(bodyJson, &mapReference)
	if err != nil {
		return "", utils.CommonError(err.Error() + " " + requestUrl + " " + string(bodyJson))
	}

	return mapReference.Url, nil
}
