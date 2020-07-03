package gpio

type GPIO interface {
	DigitalWrite(pin int, level int) (err error)
	DigitalRead(pin int) (level int, err error)
}

type Relay interface {
	On() (err error)
	Off() (err error)
	State() (state string)
}

const (
	High int = 1
	Low  int = 0
)
