package relay

const (
	on  string = "on"
	off string = "off"
)

// State represent the relay state
type State interface {
	// State return the state
	State() string

	// SetStateOn permit to set on state
	SetStateOn()

	// SetStateOff permit to set off state
	SetStateOff()

	// IsOn return true if on
	IsOn() bool

	// Is Off return true if off
	IsOff() bool
}

// StateImp implement the state interface
type StateImp struct {
	state string
}

// NewState return new state object
func NewState() (state State) {
	state = &StateImp{}
	return state
}

// State return the state
func (s *StateImp) State() string {
	return s.state
}

// SetStateOn permit to set on state
func (s *StateImp) SetStateOn() {
	s.state = on
}

// SetStateOff permit to set off state
func (s *StateImp) SetStateOff() {
	s.state = off
}

// IsOn return true if on
func (s *StateImp) IsOn() bool {
	return s.state == on
}

// Is Off return true if off
func (s *StateImp) IsOff() bool {
	return s.state == off
}
