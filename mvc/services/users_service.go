package services

import (
	"github.com/nickmonks/microservices-go/mvc/domain"
	"github.com/nickmonks/microservices-go/mvc/utils"
)

// See on user_dao.go why we return a User pointer...
func GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return domain.GetUser(userId)
}
