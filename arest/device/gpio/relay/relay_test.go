package relay

import (
	"context"
	"testing"

	"github.com/disaster37/go-arest/arest"
	"github.com/disaster37/go-arest/arest/rest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRelayRest(t *testing.T) {
	client := rest.MockRestClient()
	fixture := `{}`
	responder := httpmock.NewStringResponder(200, fixture)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/o", responder)
	httpmock.RegisterResponder("POST", "http://localhost/digital/0/1", responder)
	httpmock.RegisterResponder("POST", "http://localhost/digital/0/0", responder)

	// when NO and High signale
	signal := arest.NewLevel()
	output := NewOutput()
	defaultState := NewState()
	defaultState.SetStateOn()
	signal.SetLevelHigh()
	output.SetOutputNO()
	relay, err := NewRelay(client, 0, signal, output, defaultState)
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.On(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.Off(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOff())

	//when NO and Low signal
	signal = arest.NewLevel()
	output = NewOutput()
	defaultState = NewState()
	defaultState.SetStateOn()
	signal.SetLevelLow()
	output.SetOutputNO()
	relay, err = NewRelay(client, 0, signal, output, defaultState)
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.On(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.Off(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOff())

	// when NC and High signal
	signal = arest.NewLevel()
	output = NewOutput()
	defaultState = NewState()
	defaultState.SetStateOn()
	signal.SetLevelHigh()
	output.SetOutputNC()
	relay, err = NewRelay(client, 0, signal, output, defaultState)
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.On(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.Off(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOff())

	//when NC and Low signal
	signal = arest.NewLevel()
	output = NewOutput()
	defaultState = NewState()
	defaultState.SetStateOn()
	signal.SetLevelLow()
	output.SetOutputNC()
	relay, err = NewRelay(client, 0, signal, output, defaultState)
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.On(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOn())
	assert.Equal(t, true, relay.OutputState().IsOn())
	err = relay.Off(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, relay.State().IsOff())
	assert.Equal(t, true, relay.OutputState().IsOff())

}
