package gpio

import (
	"fmt"

	client "github.com/disaster37/go-arest"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type GPIOImp struct {
	Client client.Client
	Pin    int
	Mode   string
}

// DigitalWrite permit to set level on digital pin
func (c *GPIOImp) DigitalWrite(pin int, level int) (err error) {

	err = CheckLevel(level)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/digital/%d/%d", pin, level)

	resp, err := c.Client.Resty.R().
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	return err
}

// DigitalRead permit to read level on digital pin
func (c *GPIOImp) DigitalRead(pin int) (level int, err error) {

	url := fmt.Sprintf("/digital/%d", pin)

	resp, err := c.Client.Resty.R().
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	data, err := Unmarshal(resp.Body())
	if err != nil {
		return level, err
	}

	level = data["return_value"].(int)

	return level, err
}

// ReadValue permit to read exposed variable
func (c *Client) ReadValue(name string) (value interface{}, err error) {
	url := fmt.Sprintf("/%s", name)

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	data, err := Unmarshal(resp.Body())
	if err != nil {
		return value, err
	}

	if temp, ok := data[name]; ok {
		value = temp
	} else {
		err = errors.Errorf("Variable %s not found", name)
	}

	return value, err
}
