package led

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/disaster37/go-arest/arest/rest"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestLed(t *testing.T) {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	client := rest.MockRestClient()
	signal := arest.NewLevel()
	signal.SetLevelHigh()
	responder := httpmock.NewStringResponder(200, `{}`)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/o", responder)
	httpmock.RegisterResponder("POST", "http://localhost/digital/0/1", responder)
	httpmock.RegisterResponder("POST", "http://localhost/digital/0/0", responder)

	led, err := NewLed(client, 0, false)
	assert.NoError(t, err)

	// Turn on
	err = led.TurnOn(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)

	// Turn off
	err = led.TurnOff(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, false, led.(*LedImp).state)

	// Toogle
	led.TurnOff(context.Background())
	err = led.Toogle(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)

	// Blink
	led.TurnOff(context.Background())
	err = led.Blink(context.Background(), 1*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, false, led.(*LedImp).state)

	// Reset
	led.TurnOn(context.Background())
	err = led.Reset(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)
}
