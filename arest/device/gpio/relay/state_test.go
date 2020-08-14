package relay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {

	state := NewState()

	// On
	state.SetStateOn()
	assert.Equal(t, on, state.State())
	assert.Equal(t, true, state.IsOn())

	// Off
	state.SetStateOff()
	assert.Equal(t, off, state.State())
	assert.Equal(t, true, state.IsOff())
}
