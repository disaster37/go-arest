package arest

import "github.com/jarcoal/httpmock"

func MockClient() Arest {
	client := NewClient("http://localhost")
	httpmock.ActivateNonDefault(client.Client().GetClient())

	return client

}
