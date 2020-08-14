package arest

import "context"

// Arest permit to access on Arest API
type Arest interface {

	// SetPinMode permit to set pin mode
	SetPinMode(ctx context.Context, pin int, mode Mode) (err error)

	// DigitalWrite permit to set level on pin
	DigitalWrite(ctx context.Context, pin int, level Level) (err error)

	// DigitalRead permit to read level from pin
	DigitalRead(ctx context.Context, pin int) (level Level, err error)

	// ReadValue permit to read user variable
	ReadValue(ctx context.Context, name string) (value interface{}, err error)

	// ReadValues permit to read all user variables
	ReadValues(ctx context.Context) (values map[string]interface{}, err error)

	// CallFunction permit to call user function
	CallFunction(ctx context.Context, name string, param string) (resp int, err error)
}
