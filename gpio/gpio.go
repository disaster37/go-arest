package gpio

import client "github.com/disaster37/go-arest/v1"

type GPIO struct {
	Client client.Client
	Pin    int
	Mode   string
}
