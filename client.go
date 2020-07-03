package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	resty *resty.Client
}

func NewClient(url string) client {
	resty := resty.New().
		SetHostURL(url).
		SetHeader("Content-Type", "application/json")

	return &Client{
		resty: resty,
	}
}

// SetPinMode permit to set pin mode on target
func (c *Client) SetPinMode(pin int, mode string) (err error) {

	err = CheckMode(mode)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/mode/%d/%s", pin, mode)

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	return err

}

// DigitalWrite permit to set level on digital pin
func (c *Client) DigitalWrite(pin int, level int) (err error) {

	err = CheckLevel(level)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/digital/%d/%d", pin, level)

	resp, err := c.resty.R().
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	return err
}

// DigitalRead permit to read level on digital pin
func (c *Client) DigitalRead(pin int) (level int, err error) {

	url := fmt.Sprintf("/digital/%d", pin)

	resp, err := c.resty.R().
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

func (c *Client) CallFunction(name string, command string) (value int, err error) {

	url := fmt.Sprintf("/%s", name)

	resp, err := c.resty.R().
		SetQueryParams(map[string]string{
			"params": command,
		}).
		SetHeader("Accept", "application/json").
		Get(url)

	log.Debugf("Resp: %s", resp.String())

	data, err := Unmarshal(resp.Body())
	if err != nil {
		return value, err
	}

	if temp, ok := data["return_value"]; ok {
		value = temp.(int)
	} else {
		errors.Errorf("Function %s not found", name)
	}

	return value, err

}
