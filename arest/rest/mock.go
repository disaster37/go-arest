package rest

import (
	"github.com/disaster37/go-arest/arest"
	"github.com/jarcoal/httpmock"
)

func MockRestClient() arest.Arest {
	client := NewClient("http://localhost").(*Client)
	httpmock.ActivateNonDefault(client.Client().GetClient())

	return client

}
