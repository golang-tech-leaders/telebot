package utils

import (
	"errors"
	"log"
)

const Error404Message = "По заданным значениям данных не найдено."
const CommonErrorMessage = "Что-то пошло не так. Попробуйте снова."

func IsFailStatus(status int) bool {
	return status < 200 && status > 299
}

func HttpError(status int, info string) error {
	log.Println("ERROR " + info)
	if status == 404 {
		return errors.New(Error404Message)
	}
	return errors.New(CommonErrorMessage)
}

func HttpErrorWithCustom404(status int, info string, message404 string) error {
	log.Println("ERROR " + info)
	if status == 404 {
		return errors.New(message404)
	}
	return errors.New(CommonErrorMessage)
}

func CommonError(info string) error {
	log.Println("ERROR " + info)
	return errors.New(CommonErrorMessage)
}
