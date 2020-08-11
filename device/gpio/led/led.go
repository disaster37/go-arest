package led

import (
	"time"

	"github.com/disaster37/go-arest"
	log "github.com/sirupsen/logrus"
)

// Led is the led interface
type Led interface {
	TurnOn() error
	TurnOff() error
	Blink(duration time.Duration) *time.Timer
	Toogle() error
	Reset() error
}

// LedImp is the default Led implementation
type LedImp struct {
	pin    int
	client arest.Arest
	state  bool
}

// NewLed return new led device
func NewLed(client arest.Arest, pin int, defaultState bool) (Led, error) {

	led := &LedImp{
		pin:    pin,
		client: client,
		state:  defaultState,
	}

	err := led.Reset()
	if err != nil {
		return nil, err
	}

	return led, nil
}

// TurnOn turn on led
func (h *LedImp) TurnOn() error {
	level := arest.NewLevel()
	level.SetLevelHigh()
	err := h.client.DigitalWrite(h.pin, level)
	if err != nil {
		return err
	}

	h.state = true
	return nil
}

// TurnOff turn on led
func (h *LedImp) TurnOff() error {
	level := arest.NewLevel()
	level.SetLevelLow()
	err := h.client.DigitalWrite(h.pin, level)
	if err != nil {
		return err
	}

	h.state = false
	return nil
}

// Toogle the led state
func (h *LedImp) Toogle() error {
	level := arest.NewLevel()
	if h.state {
		level.SetLevelLow()
	} else {
		level.SetLevelHigh()
	}
	err := h.client.DigitalWrite(h.pin, level)
	if err != nil {
		return err
	}

	h.state = level.IsHigh()
	return nil
}

// Reset put the led on desired state
func (h *LedImp) Reset() error {

	mode := arest.NewMode()
	mode.SetModeOutput()
	err := h.client.SetPinMode(h.pin, mode)
	if err != nil {
		return err
	}

	if h.state {
		return h.TurnOn()
	}

	return h.TurnOff()
}

// Blink the led during time
func (h *LedImp) Blink(duration time.Duration) *time.Timer {

	timer := time.NewTimer(duration)
	quit := make(chan bool)

	// Lauch timer for blink duration duration
	go func() {
		<-timer.C
		quit <- true
	}()

	// Start a loop to blink led
	go func() {
		expectedState := h.state
		for {
			select {
			case <-quit:
				if expectedState {
					err := h.TurnOn()
					if err != nil {
						log.Errorf("Error appear when turn on led: %s", err.Error())
					}
				} else {
					err := h.TurnOff()
					if err != nil {
						log.Errorf("Error appear when turn off led: %s", err.Error())
					}
				}
				return
			default:
				err := h.Toogle()
				if err != nil {
					log.Errorf("Error appear when toogle led: %s", err.Error())
				}
				time.Sleep(1 * time.Second)
			}
		}

	}()

	return timer
}
