package modules

type User struct {
	Name     string
	Email    string
	Password string
}

type UserWithoutPassword struct {
	Name  string
	Email string
}
