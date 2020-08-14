package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var errContextAborded error = errors.Errorf("Aborded by context")

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
func (c *Client) SetPinMode(ctx context.Context, pin int, mode arest.Mode) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Debugf("Pin: %d, Mode: %s", pin, mode.String())

		url := fmt.Sprintf("/mode/%d/%s", pin, mode.Mode())

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			Post(url)

		log.Debugf("Resp: %s", resp.String())

		return err
	}

	return err

}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(ctx context.Context, pin int, level arest.Level) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Debugf("Pin: %d, Level: %s", pin, level.String())

		url := fmt.Sprintf("/digital/%d/%d", pin, level.Level())

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			Post(url)

		log.Debugf("Resp: %s", resp.String())

		return err
	}

	return err

}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(ctx context.Context, pin int) (level arest.Level, err error) {

	select {
	case <-ctx.Done():
		return level, ctx.Err()
	default:

		log.Debugf("Pin: %d", pin)

		url := fmt.Sprintf("/digital/%d", pin)
		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			SetResult(&data).
			Get(url)
		if err != nil {
			return nil, err
		}

		log.Debugf("Resp: %s, %+v", resp.String(), data)

		level = arest.NewLevel()
		if int(data["return_value"].(float64)) == arest.High {
			level.SetLevelHigh()
		} else {
			level.SetLevelLow()
		}

		return level, err
	}

	return level, err
}

// ReadValue permit to read user variable
func (c *Client) ReadValue(ctx context.Context, name string) (value interface{}, err error) {

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	default:

		log.Debugf("Value name: %s", name)

		url := fmt.Sprintf("/%s", name)
		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
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

	return value, err
}

// ReadValues permit to read user variable
func (c *Client) ReadValues(ctx context.Context) (values map[string]interface{}, err error) {
	select {
	case <-ctx.Done():
		return values, ctx.Err()
	default:

		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
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

	return values, err
}

// CallFunction permit to call user function
func (c *Client) CallFunction(ctx context.Context, name string, param string) (value int, err error) {

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	default:

		log.Debugf("Function: %s, param: %s", name, param)

		url := fmt.Sprintf("/%s", name)

		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetQueryParams(map[string]string{
				"params": param,
			}).
			SetHeader("Accept", "application/json").
			SetContext(ctx).
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

	return value, err

}
