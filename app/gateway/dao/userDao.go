package dao

import (
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}
