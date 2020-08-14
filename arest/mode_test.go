package arest

import "github.com/stretchr/testify/assert"

func (s *ArestTestSuite) testMode() {

	mode := NewMode()

	// output
	mode.SetModeOutput()
	assert.Equal(s.T(), output, mode.Mode())
	assert.Equal(s.T(), "output", mode.String())

	// input
	mode.SetModeInput()
	assert.Equal(s.T(), input, mode.Mode())
	assert.Equal(s.T(), "input", mode.String())

	// input_pullup
	mode.SetModeInputPullup()
	assert.Equal(s.T(), input_pullup, mode.Mode())
	assert.Equal(s.T(), "input_pullup", mode.String())
}
