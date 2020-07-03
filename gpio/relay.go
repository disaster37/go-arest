package gpio

import (
	"errors"

	client "github.com/disaster37/go-arest"
)

const (
	NO       int = 0
	NC       int = 1
	StateOn      = "on"
	StateOff     = "off"
)

type RelayImp struct {
	GPIO         *GPIO
	Level        int
	Output       int
	DefaultState string
	state        string
}

func NewRelay(c client.Client, pin int, level int, output int, defaultState string) (relay Relay, err error) {

	if level != High && level != Low {
		errors.New("Level must be High or Low")
	}
	if output != NO && output != NC {
		errors.New("Output must be NO or NC")
	}
	if defaultState != StateOn && defaultState != StateOff {
		errors.New("DefaultState must be StateOn or StateOff")
	}

	gpio := &GPIOImp{
		Client: c,
		Pin:    pin,
		Mode:   client.OUTPUT,
	}

	relay = &RelayImp{
		GPIO:         gpio,
		Level:        level,
		Output:       output,
		DefaultState: defaultState,
	}

	// Set pin mode
	err = c.SetPinMode(pin, gpio.Mode)
	if err != nil {
		return nil, err
	}

	// Set default state
	switch defaultState {
	case StateOn:
		err = relay.On()
	case StateOff:
		err = relay.Off()
	}

	return relay, err

}

func (r *RelayImp) On() (err error) {

	switch r.Output {
	case NO:
		// Normaly Open
		switch r.Level {
		case High:
			// High signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.HIGH)
		case Low:
			// Low signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.LOW)
		}
	case NC:
		// Normaly Close
		switch r.Level {
		case High:
			// High signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.LOW)
		case Low:
			// Low signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.HIGH)
		}
	}

	if err == nil {
		r.state = StateOn
	}

	return err
}

func (r *RelayImp) Off() (err error) {
	switch r.Output {
	case NO:
		// Normaly Open
		switch r.Level {
		case High:
			// High signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.LOW)
		case Low:
			// Low signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.HIGH)
		}
	case NC:
		// Normaly Close
		switch r.Level {
		case High:
			// High signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.HIGH)
		case Low:
			// Low signal
			err = r.GPIO.DigitalWrite(r.GPIO.Pin, client.LOW)
		}
	}

	if err == nil {
		r.state = StateOff
	}

	return err
}

func (r *RelayImp) State() string {
	return r.State()
}
