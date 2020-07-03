package gpio

type Relay interface {
	On() (err error)
	Off() (err error)
	State() (state string)
}
