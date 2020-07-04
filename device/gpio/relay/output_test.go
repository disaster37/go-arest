package relay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {

	output := NewOutput()

	// NO
	output.SetOutputNO()
	assert.Equal(t, no, output.Output())
	assert.Equal(t, true, output.IsNO())

	// NC
	output.SetOutputNC()
	assert.Equal(t, nc, output.Output())
	assert.Equal(t, true, output.IsNC())
}