package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"telebot/internal/http_error"
	"telebot/internal/models"
)

// GetGeoUrl Получение ссылки из geobase
func GetGeoUrl(url string, wasteTypeId int, lat float64, lon float64) (*string, error) {
	requestUrl := fmt.Sprintf("%swaste/type/%d/location?latitude=%f&longitude=%f", url, wasteTypeId, lat, lon)
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
		return nil, http_error.HttpError(resp.StatusCode, resp.Status+" "+requestUrl+" "+string(bodyJson))
	}

	var mapReference models.MapReference
	err = json.Unmarshal(bodyJson, &mapReference)
	if err != nil {
		return nil, http_error.CommonError(err.Error() + " " + requestUrl + " " + string(bodyJson))
	}

	return &mapReference.Url, nil
}
