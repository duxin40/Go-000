package dao

import (
	"database/sql"
	"pkg/mod/github.com/pkg/errors@v0.9.1"
)

type UserInfo struct {
	ID 		int64  `json:"id"`
	Name 	string `json:"name"`
}

func (u *UserInfo) GetName() (string, error){
	err := sql.ErrNoRows
	return "", errors.Wrap(err, "dao: get user name from db fail")
}
