package service

import (
	"gogo/app/gateway/dao"
	"gogo/app/gateway/model"
)

type UserService struct {
	dao.UserDao
}

func (us UserService) InsertUser(u model.User) int {
	return u.Id
}
