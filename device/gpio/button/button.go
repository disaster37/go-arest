package button

import "github.com/disaster37/go-arest"

// Button is the button interface
type Button interface {
	IsPushed() (state bool)
	IsReleazed() (state bool)
	IsUp() (state bool)
	IsDown() (state bool)
	Read() error
}

// ButtonImp is the default Button implementation
type ButtonImp struct {
	pin         int
	client      arest.Arest
	signal      arest.Level
	inputPullup bool
	isPushed    bool
	isReleazed  bool
	state       bool
}

func NewButton(client arest.Arest, pin int, signal arest.Level, inputPullup bool) (Button, error) {
	mode := arest.NewMode()
	if !inputPullup {
		mode.SetModeInput()
	} else {
		mode.SetModeInputPullup()
	}

	err := client.SetPinMode(pin, mode)
	if err != nil {
		return nil, err
	}

	return &ButtonImp{
		pin:         pin,
		client:      client,
		signal:      signal,
		inputPullup: inputPullup,
		isPushed:    false,
		isReleazed:  false,
		state:       false,
	}, nil
}

// IsPushed return true if button just to be pushed
func (h *ButtonImp) IsPushed() (state bool) {
	return h.isPushed
}

// IsReleazed return true if button jus to be releazed
func (h *ButtonImp) IsReleazed() (state bool) {
	return h.isReleazed
}

// IsUp return true if button  is UP
func (h *ButtonImp) IsUp() (state bool) {
	return !h.state
}

// IsDown return true if button is down
func (h *ButtonImp) IsDown() (state bool) {
	return h.state
}

func (h *ButtonImp) Read() error {
	level, err := h.client.DigitalRead(h.pin)
	if err != nil {
		return err
	}

	var computedLevel bool
	if h.signal.IsHigh() {
		computedLevel = level.IsHigh()
	} else {
		computedLevel = level.IsLow()
	}
	if h.inputPullup {
		computedLevel = !computedLevel
	}

	if computedLevel && !h.state {
		// Just push button
		h.isPushed = true
		h.isReleazed = false
	} else if !computedLevel && h.state {
		// Juste releaze button
		h.isPushed = false
		h.isReleazed = true
	} else {
		h.isPushed = false
		h.isReleazed = false
	}
	h.state = computedLevel

	return nil
}
