package button

import (
	"testing"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/rest"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestInputButton(t *testing.T) {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	client := rest.MockRestClient()
	signal := arest.NewLevel()
	signal.SetLevelHigh()
	responderMode := httpmock.NewStringResponder(200, `{}`)
	responderUp := httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
		"return_value": 0,
	})
	responderDown := httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
		"return_value": 1,
	})
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/i", responderMode)

	button, err := NewButton(client, 0, signal)
	assert.NoError(t, err)

	// When read button on default state
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When push button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderDown)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, button.IsDown())
	assert.Equal(t, false, button.IsUp())
	assert.Equal(t, true, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When keep button pushed
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderDown)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, button.IsDown())
	assert.Equal(t, false, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When releaze button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, true, button.IsReleazed())

	// When keep releaze button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

}

func TestInputPullupButton(t *testing.T) {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	client := rest.MockRestClient()
	signal := arest.NewLevel()
	signal.SetLevelLow()
	responderMode := httpmock.NewStringResponder(200, `{}`)
	responderUp := httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
		"return_value": 1,
	})
	responderDown := httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
		"return_value": 0,
	})
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/I", responderMode)

	button, err := NewButton(client, 0, signal)
	assert.NoError(t, err)

	// When read button on default state
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When push button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderDown)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, button.IsDown())
	assert.Equal(t, false, button.IsUp())
	assert.Equal(t, true, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When keep button pushed
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderDown)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, button.IsDown())
	assert.Equal(t, false, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

	// When releaze button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, true, button.IsReleazed())

	// When keep releaze button
	httpmock.Reset()
	httpmock.RegisterResponder("GET", "http://localhost/digital/0", responderUp)
	err = button.Read()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, button.IsDown())
	assert.Equal(t, true, button.IsUp())
	assert.Equal(t, false, button.IsPushed())
	assert.Equal(t, false, button.IsReleazed())

}
