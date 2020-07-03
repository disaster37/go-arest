package client

type Board interface {
	SetPinMode(pin int, mode string) (err error)
	ReadValue(name string) (value interface{}, err error)
	CallFunction(name string, command string) (resp int, err error)
}

const (
	INPUT        string = "i"
	INPUT_PULLUP string = "I"
	OUTPUT       string = "o"
)
