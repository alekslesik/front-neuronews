package model

type User struct {
	ID       int
	Username string
	Email    string
}

type UserModel interface {
	GetUserByID(id int) (*User, error)
}
