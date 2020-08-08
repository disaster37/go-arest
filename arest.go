package arest

// Arest permit to access on Arest API
type Arest interface {

	// SetPinMode permit to set pin mode
	SetPinMode(pin int, mode Mode) (err error)

	// DigitalWrite permit to set level on pin
	DigitalWrite(pin int, level Level) (err error)

	// DigitalRead permit to read level from pin
	DigitalRead(pin int) (level Level, err error)

	// ReadValue permit to read user variable
	ReadValue(name string) (value interface{}, err error)

	// CallFunction permit to call user function
	CallFunction(name string, param string) (resp int, err error)
}
