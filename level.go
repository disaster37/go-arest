package arest

const (
	High int = 1
	Low  int = 0
)

// Level is the digital pin level (HIGH or LOW)
type Level interface {
	// Level return the current pin level
	Level() int

	// SetLevelHigh permit to set the high level
	SetLevelHigh()

	// SetLevelLow permit to set the low level
	SetLevelLow()

	// String return the current level as human name
	String() string

	// IsHigh return true if level is high
	IsHigh() bool

	// IsLow return true if level is low
	IsLow() bool
}

// LevelImp is the implementation of Level interface
type LevelImp struct {
	level int
}

// NewLevel return new level Obj
func NewLevel() (level Level) {
	level = &LevelImp{}
	return level
}

// Level return the current pin level
func (l *LevelImp) Level() int {
	return l.level
}

// SetLevelHigh permit to set the high level
func (l *LevelImp) SetLevelHigh() {
	l.level = High
}

// SetLevelLow permit to set the high level
func (l *LevelImp) SetLevelLow() {
	l.level = Low
}

// String return the current level as human name
func (l *LevelImp) String() string {
	switch l.level {
	case High:
		return "high"
	case Low:
		return "low"
	}

	return ""
}

// IsHigh return true if level is high
func (l *LevelImp) IsHigh() bool {
	return l.level == High
}

// IsLow return true if level is low
func (l *LevelImp) IsLow() bool {
	return l.level == Low
}
