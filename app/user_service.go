package app

// UserService ..
type UserService interface {
	Login() (interface{}, string, error)
	Logout() (interface{}, string, error)
	Authenticate(token string) (interface{}, string, error)
}
