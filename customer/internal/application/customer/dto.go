package customer

type RegisterDto struct {
	Name     string
	Phone    string
	Password string
}

type LoginDto struct {
	Phone    string
	Password string
}
