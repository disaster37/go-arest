package gpio

import "errors"

// CheckLevel permit to check that level provided exist
func CheckLevel(level int) (err error) {
	if level != High && level != Low {
		err = errors.New("Level must be High or Low")
	}

	return err
}
