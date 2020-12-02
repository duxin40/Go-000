package controller

import (
	"Go-000/Week02/service"
	"log"
	"pkg/mod/github.com/pkg/errors@v0.9.1"
)

func GetUsernameByIDController(id int64) (string, error) {
	name, err := service.GetUserNameByID(id)
	if err != nil {
		log.Printf("get user name by id fail, error with stack: %+v", err)
		return "", errors.WithMessage(err, "controller: get username by id fail")
	}
	return name , nil
}
