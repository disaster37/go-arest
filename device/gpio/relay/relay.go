package relay

import (
	"github.com/disaster37/go-arest"
)

// Relay represent relay device
type Relay interface {

	// On enable the relay output
	On() (err error)

	// Off disable the relay output
	Off() (err error)

	// State return the current relay state
	State() (state State)

	// OutputState return the current output state
	OutputState() (state State)
}

// RelayImp implement the relay interface
type RelayImp struct {
	pin         int
	client      arest.Arest
	signal      arest.Level
	output      Output
	state       State
	outputState State
}

// NewRelay return new relay object
func NewRelay(c arest.Arest, pin int, signal arest.Level, output Output, defaultState State) (relay Relay, err error) {

	relay = &RelayImp{
		client:      c,
		pin:         pin,
		signal:      signal,
		output:      output,
		state:       NewState(),
		outputState: NewState(),
	}

	mode := arest.NewMode()
	mode.SetModeOutput()

	// Set pin mode
	err = c.SetPinMode(pin, mode)
	if err != nil {
		return nil, err
	}

	// Set default state
	if defaultState.IsOn() {
		err = relay.On()
	} else {
		err = relay.Off()
	}

	return relay, err

}

// On enable the relay output
func (r *RelayImp) On() (err error) {

	level := arest.NewLevel()
	state := NewState()

	if r.output.IsNO() {
		if r.signal.IsHigh() {
			level.SetLevelHigh()
			state.SetStateOn()
		} else {
			level.SetLevelLow()
			state.SetStateOff()
		}
	} else {
		if r.signal.IsHigh() {
			level.SetLevelLow()
			state.SetStateOff()
		} else {
			level.SetLevelHigh()
			state.SetStateOn()
		}
	}

	err = r.client.DigitalWrite(r.pin, level)
	if err != nil {
		return err
	}

	r.outputState.SetStateOn()
	if state.IsOn() {
		r.state.SetStateOn()
	} else {
		r.state.SetStateOff()
	}

	return nil
}

// Off disable the relay output
func (r *RelayImp) Off() (err error) {
	level := arest.NewLevel()
	state := NewState()

	if r.output.IsNO() {
		if r.signal.IsHigh() {
			level.SetLevelLow()
			state.SetStateOff()
		} else {
			level.SetLevelHigh()
			state.SetStateOn()
		}
	} else {
		if r.signal.IsHigh() {
			level.SetLevelHigh()
			state.SetStateOn()
		} else {
			level.SetLevelLow()
			state.SetStateOff()
		}
	}

	err = r.client.DigitalWrite(r.pin, level)
	if err != nil {
		return err
	}

	r.outputState.SetStateOff()
	if state.IsOn() {
		r.state.SetStateOn()
	} else {
		r.state.SetStateOff()
	}

	return nil
}

// State return the current relay state
func (r *RelayImp) State() State {
	return r.state
}

// OutputState return the current output state
func (r *RelayImp) OutputState() State {
	return r.outputState
}
