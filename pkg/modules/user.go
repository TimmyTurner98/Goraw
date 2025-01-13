package modules

type User struct {
	Id       int `json:"id"`
	Name     string
	Email    string
	Password string
}

type UserWithoutPassword struct {
	Id    int `json:"id"`
	Name  string
	Email string
}
