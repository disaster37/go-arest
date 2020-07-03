package gpio

import "github.com/disaster37/go-client-arest/v1"

type GPIO struct {
	Client client.Client
	Pin    int
	Mode   string
}
