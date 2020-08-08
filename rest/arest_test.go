package rest

import (
	"testing"

	"github.com/disaster37/go-arest"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type ArestTestSuite struct {
	suite.Suite
	client arest.Arest
}

func (s *ArestTestSuite) SetupSuite() {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)
}

func (s *ArestTestSuite) BeforeSuite() {
	s.client = MockRestClient()
}

func (s *ArestTestSuite) AfterSuite() {
	httpmock.DeactivateAndReset()
}

func (s *ArestTestSuite) BeforeEach() {
	httpmock.Reset()
}

func TestArestTestSuite(t *testing.T) {
	suite.Run(t, new(ArestTestSuite))
}

func (s *ArestTestSuite) testSetMode() {

	fixture := `{"message": "Pin D0 set to output", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/mode/0/o"
	httpmock.RegisterResponder("POST", fakeURL, responder)

	mode := arest.NewMode()
	mode.SetModeOutput()
	err := s.client.SetPinMode(0, mode)
	assert.NoError(s.T(), err)

}

func (s *ArestTestSuite) testDigitalWrite() {

	fixture := `{"message": "Pin D0 set to 1", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/digital/0/1"
	httpmock.RegisterResponder("POST", fakeURL, responder)

	level := arest.NewLevel()
	level.SetLevelHigh()
	err := s.client.DigitalWrite(0, level)
	assert.NoError(s.T(), err)
}

func (s *ArestTestSuite) testDigitalRead() {

	fixture := `{"return_value": 1, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/digital/0"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	level, err := s.client.DigitalRead(0)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "high", level.String())
}

func (s *ArestTestSuite) testReadValue() {

	fixture := `{"isRebooted": true, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/isRebooted"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	value, err := s.client.ReadValue("isRebooted")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, value.(bool))

	// Bad
	value, err = s.client.ReadValue("bad")
	assert.Error(s.T(), err)
}

func (s *ArestTestSuite) testCallFunction() {

	fixture := `{"return_value": 1, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/acknoledgeRebooted?params=test"
	httpmock.RegisterResponder("POST", fakeURL, responder)

	resp, err := s.client.CallFunction("acknoledgeRebooted", "test")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, resp)

	// Bad
	resp, err = s.client.CallFunction("bad", "test")
	assert.Error(s.T(), err)
}
