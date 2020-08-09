package led

import (
	"testing"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/rest"
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
	err = led.TurnOn()
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)

	// Turn off
	err = led.TurnOff()
	assert.NoError(t, err)
	assert.Equal(t, false, led.(*LedImp).state)

	// Toogle
	led.TurnOff()
	err = led.Toogle()
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)

	// Blink
	led.TurnOff()
	timer := led.Blink(1 * time.Second)
	<-timer.C
	time.Sleep(5 * time.Second)
	assert.Equal(t, false, led.(*LedImp).state)

	// Reset
	led.TurnOn()
	err = led.Reset()
	assert.NoError(t, err)
	assert.Equal(t, true, led.(*LedImp).state)
}
