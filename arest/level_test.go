package arest

import "github.com/stretchr/testify/assert"

func (s *ArestTestSuite) testLevel() {

	level := NewLevel()

	// High
	level.SetLevelHigh()
	assert.Equal(s.T(), High, level.Level())
	assert.Equal(s.T(), "high", level.String())
	assert.Equal(s.T(), true, level.IsHigh())

	// Low
	level.SetLevelLow()
	assert.Equal(s.T(), Low, level.Level())
	assert.Equal(s.T(), "low", level.String())
	assert.Equal(s.T(), true, level.IsLow())
}
