// IMPORTANT: domain is equivalent to Model in the MVC pattern
package domain

type User struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
