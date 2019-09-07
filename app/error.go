package app

// SingleError ..
type SingleError struct {
	Slug    string
	Message string
	details struct{ path string }
}

// Error ..
type Error struct {
	Errors  []SingleError `json:"error"`
	Message string        `json:"message"`
}

// NewError ..
func NewError(msg string) Error {
	errDetails := struct{ path string }{"Input"}
	singleError := SingleError{"Invlaide-Input", msg, errDetails}
	return Error{[]SingleError{singleError}, msg}
}
