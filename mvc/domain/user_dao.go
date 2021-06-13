package domain

import (
	"fmt"
	"net/http"

	"github.com/nickmonks/microservices-go/mvc/utils"
)

var (
	users = map[int64]*User{
		123: {Id: 123, FirstName: "Fede", LastName: "Monkerud", Email: "nicolas@gmail.com"},
	}
)

// Access to database
// IMPORTANT: if we return a *User instead of User, we dont need to return an empty User{}, we can simply return nil
func GetUser(userId int64) (*User, *utils.ApplicationError) {
	user := users[userId]
	if user == nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("user %v was not found", userId),
			StatusCode: http.StatusNotFound,
			Code:       "Not found",
		}
	}

	return user, nil
}
