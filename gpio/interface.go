package gpio

type relay interface {
	On() (err error)
	Off() (err error)
	State() (state string)
}
