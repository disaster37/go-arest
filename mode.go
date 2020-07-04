package arest

const (
	input        string = "i"
	input_pullup string = "I"
	output       string = "o"
)

// Mode permit to set pin mode
type Mode interface {
	// Mode get the current mode
	Mode() string

	// SetModeOutput configure mode as output
	SetModeOutput()

	// SetModeInput configure mode as input
	SetModeInput()

	// SetModeInputPullup configure mode as inputPullup
	SetModeInputPullup()

	// String get the human mode
	String() string
}

// ModeImp implement mode interface
type ModeImp struct {
	mode string
}

// NewMode initialize new mode object
func NewMode() (mode Mode) {
	mode = &ModeImp{}

	return mode
}

// Mode permit to get the current mode
func (m *ModeImp) Mode() string {
	return m.mode
}

// SetModeInput permit to set the input mode
func (m *ModeImp) SetModeInput() {
	m.mode = input
}

// SetModeInput permit to set the input mode
func (m *ModeImp) SetModeInputPullup() {
	m.mode = input_pullup
}

// SetModeInput permit to set the input mode
func (m *ModeImp) SetModeOutput() {
	m.mode = output
}

// String get the human mode
func (m *ModeImp) String() string {
	switch m.mode {
	case input:
		return "input"
	case input_pullup:
		return "input_pullup"
	case output:
		return "output"
	}

	return ""
}
