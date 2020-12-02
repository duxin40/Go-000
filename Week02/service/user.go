package service

import (
	"Go-000/Week02/dao"
	"pkg/mod/github.com/pkg/errors@v0.9.1"
)

func GetUserNameByID(id int64) (string, error) {
	u := &dao.UserInfo{
		ID:   id,
	}
	name, err := u.GetName()
	if err != nil {
		return "", errors.WithMessage(err, "service: get username by id fail")
	}
	return name , nil
}
