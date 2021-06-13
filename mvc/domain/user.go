// IMPORTANT: domain is equivalent to Model in the MVC pattern
package domain

type User struct {
	Id        uint64
	FirstName string
	LastName  string
	Email     string
}
