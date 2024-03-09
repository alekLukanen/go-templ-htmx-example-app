package database

type User struct {
	id   int
	name string
}

func GetUser(id int) User {
	return User{id: id}
}
