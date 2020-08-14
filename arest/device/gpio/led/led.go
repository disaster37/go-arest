package led

import (
	"context"
	"time"

	"github.com/disaster37/go-arest/arest"
)

// Led is the led interface
type Led interface {
	TurnOn(ctx context.Context) error
	TurnOff(ctx context.Context) error
	Blink(ctx context.Context, duration time.Duration) error
	Toogle(ctx context.Context) error
	Reset(ctx context.Context) error
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

	err := led.Reset(context.Background())
	if err != nil {
		return nil, err
	}

	return led, nil
}

// TurnOn turn on led
func (h *LedImp) TurnOn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		level := arest.NewLevel()
		level.SetLevelHigh()
		err := h.client.DigitalWrite(ctx, h.pin, level)
		if err != nil {
			return err
		}

		h.state = true
		return nil
	}

	return nil

}

// TurnOff turn on led
func (h *LedImp) TurnOff(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		level := arest.NewLevel()
		level.SetLevelLow()
		err := h.client.DigitalWrite(ctx, h.pin, level)
		if err != nil {
			return err
		}

		h.state = false
		return nil
	}

	return nil

}

// Toogle the led state
func (h *LedImp) Toogle(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		level := arest.NewLevel()
		if h.state {
			level.SetLevelLow()
		} else {
			level.SetLevelHigh()
		}
		err := h.client.DigitalWrite(ctx, h.pin, level)
		if err != nil {
			return err
		}

		h.state = level.IsHigh()
		return nil
	}

	return nil

}

// Reset put the led on desired state
func (h *LedImp) Reset(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		mode := arest.NewMode()
		mode.SetModeOutput()
		err := h.client.SetPinMode(ctx, h.pin, mode)
		if err != nil {
			return err
		}

		if h.state {
			return h.TurnOn(ctx)
		}

		return h.TurnOff(ctx)
	}

	return nil

}

// Blink the led during time
func (h *LedImp) Blink(ctx context.Context, duration time.Duration) error {

	timer := time.NewTimer(duration)
	ch := make(chan bool)
	chErr := make(chan error)
	expectedState := h.state

	go func() {
		for {
			select {
			case <-ch:
				if expectedState {
					err := h.TurnOn(ctx)
					if err != nil {
						chErr <- err
					}
				} else {
					err := h.TurnOff(ctx)
					if err != nil {
						chErr <- err
					}
				}

				chErr <- nil
				return
			default:
				err := h.Toogle(ctx)
				if err != nil {
					chErr <- err
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	select {
	case <-ctx.Done():
		ch <- true
		err := <-chErr
		if err != nil {
			return err
		}
		return ctx.Err()
	case <-timer.C:
		ch <- true
		err := <-chErr
		if err != nil {
			return err
		}
	case err := <-chErr:
		return err
	}

	return nil
}
