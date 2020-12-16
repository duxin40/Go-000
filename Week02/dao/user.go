package dao

import (
	"database/sql"
	"github.com/pkg/errors"
)

type UserInfo struct {
	ID 		int64  `json:"id"`
	Name 	string `json:"name"`
}

func (u *UserInfo) GetName() (string, error){
	err := sql.ErrNoRows
	return "", errors.Wrap(err, "dao: get user name from db fail")
}
