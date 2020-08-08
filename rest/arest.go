package rest

import (
	"fmt"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client implement arest interface
type Client struct {
	resty *resty.Client
}

// NewClient permit to initialize new client Object
func NewClient(url string) arest.Arest {
	resty := resty.New().
		SetHostURL(url).
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second)

	return &Client{
		resty: resty,
	}
}

// Client permit to get curent resty client
func (c *Client) Client() *resty.Client {
	return c.resty
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(pin int, mode arest.Mode) (err error) {

	log.Debugf("Pin: %d, Mode: %s", pin, mode.String())

	url := fmt.Sprintf("/mode/%d/%s", pin, mode.Mode())

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		Post(url)

	log.Debugf("Resp: %s", resp.String())

	return err

}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(pin int, level arest.Level) (err error) {

	log.Debugf("Pin: %d, Level: %s", pin, level.String())

	url := fmt.Sprintf("/digital/%d/%d", pin, level.Level())

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		Post(url)

	log.Debugf("Resp: %s", resp.String())

	return err
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(pin int) (level arest.Level, err error) {

	log.Debugf("Pin: %d", pin)

	url := fmt.Sprintf("/digital/%d", pin)
	data := make(map[string]interface{})

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		SetResult(&data).
		Get(url)
	if err != nil {
		return nil, err
	}

	log.Debugf("Resp: %s", resp.String())

	level = arest.NewLevel()
	if int(data["return_value"].(float64)) == arest.High {
		level.SetLevelHigh()
	} else {
		level.SetLevelLow()
	}

	return level, err
}

// ReadValue permit to read user variable
func (c *Client) ReadValue(name string) (value interface{}, err error) {

	log.Debugf("Value name: %s", name)

	url := fmt.Sprintf("/%s", name)
	data := make(map[string]interface{})

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		SetResult(&data).
		Get(url)
	if err != nil {
		return nil, err
	}

	log.Debugf("Resp: %s", resp.String())

	if temp, ok := data[name]; ok {
		value = temp
	} else {
		err = errors.Errorf("Variable %s not found", name)
	}

	return value, err
}

// ReadValues permit to read user variable
func (c *Client) ReadValues() (values map[string]interface{}, err error) {

	data := make(map[string]interface{})

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		SetResult(&data).
		Get("/")
	if err != nil {
		return nil, err
	}

	log.Debugf("Resp: %s", resp.String())

	if temp, ok := data["variables"]; ok {
		values = temp.(map[string]interface{})
	} else {
		err = errors.Errorf("No variable found")
	}

	return values, err
}

// CallFunction permit to call user function
func (c *Client) CallFunction(name string, param string) (value int, err error) {

	log.Debugf("Function: %s, param: %s", name, param)

	url := fmt.Sprintf("/%s", name)

	data := make(map[string]interface{})

	resp, err := c.resty.R().
		SetQueryParams(map[string]string{
			"params": param,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&data).
		Post(url)
	if err != nil {
		return value, err
	}

	log.Debugf("Resp: %s", resp.String())

	if temp, ok := data["return_value"]; ok {
		value = int(temp.(float64))
	} else {
		errors.Errorf("Function %s not found", name)
	}

	return value, err

}
