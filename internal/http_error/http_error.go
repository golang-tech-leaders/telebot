package http_error

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrNotFound = errors.New("По заданным значениям данных не найдено")
	ErrUnknown  = errors.New("Что-то пошло не так. Попробуйте снова")
)

func IsFailStatus(status int) bool {
	return status < http.StatusOK || status > http.StatusIMUsed
}

func HttpError(status int, info string) error {
	log.Println("ERROR " + info)
	if status == http.StatusNotFound {
		return ErrNotFound
	}
	return ErrUnknown
}

func HttpErrorWithCustom404(status int, info string, message404 string) error {
	log.Println("ERROR " + info)
	if status == http.StatusNotFound {
		return errors.New(message404)
	}
	return ErrUnknown
}

func CommonError(info string) error {
	log.Println("ERROR " + info)
	return ErrUnknown
}
