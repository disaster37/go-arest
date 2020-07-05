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

	// Reset permit to reconfigure relay. It usefull when board reboot
	Reset() (err error)
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
		outputState: defaultState,
	}

	err = relay.Reset()

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

// Reset permit to reconfigure relay. It usefull when board reboot
// It apply the desired state
func (r *RelayImp) Reset() (err error) {
	mode := arest.NewMode()
	mode.SetModeOutput()

	// Set pin mode
	err = r.client.SetPinMode(r.pin, mode)
	if err != nil {
		return err
	}

	// Set relay on right state
	if r.outputState.IsOn() {
		err = r.On()
	} else {
		err = r.Off()
	}

	return err

}
