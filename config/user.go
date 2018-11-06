package config

func newUser() *User {
	return &User{}
}

func newUserWithDefaults() *User {
	return &User{
		Name:  "",
		Email: "",
	}
}

type User struct {
	Name  string
	Email string
}
