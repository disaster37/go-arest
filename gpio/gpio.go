package gpio

import "github.com/disaster37/go-client-arest"

type GPIO struct {
	Client client.Client
	Pin    int
	Mode   string
}
