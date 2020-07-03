package client

type Board interface {
	SetPinMode(pin int, mode string) (err error)
	DigitalWrite(pin int, level int) (err error)
	DigitalRead(pin int) (level int, err error)
	ReadValue(name string) (value interface{}, err error)
	CallFunction(name string, command string) (resp int, err error)
}

const (
	HIGH         int    = 1
	LOW          int    = 0
	INPUT        string = "i"
	INPUT_PULLUP string = "I"
	OUTPUT       string = "o"
)
