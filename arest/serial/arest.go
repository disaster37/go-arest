package serial

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/pkg/errors"
	"go.bug.st/serial"
)

// Client implement arest interface
type Client struct {
	serialPort serial.Port
	sem        chan int
	timeout    time.Duration
	url        string
}

// NewClient permit to initialize new client Object
func NewClient(url string, timeout time.Duration, debug bool) (arest.Arest, error) {

	if debug {
		arest.IsDebug = true
	}

	serialPort, err := open(url)
	if err != nil {
		return nil, err
	}

	client := &Client{
		serialPort: serialPort,
		sem:        make(chan int, 1),
		timeout:    timeout,
		url:        url,
	}

	return client, nil
}

// Client permit to get curent resty client
func (c *Client) Client() serial.Port {
	return c.serialPort
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(ctx context.Context, pin int, mode arest.Mode) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.takeSemaphore()
		defer c.releazeSemaphore()
		arest.Debug("Pin: %d, Mode: %s", pin, mode.String())

		url := fmt.Sprintf("/mode/%d/%s\n\r", pin, mode.Mode())

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return err
		}

		var resp string
		select {
		case err := <-chErr:
			return err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		return nil
	}

	return err

}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(ctx context.Context, pin int, level arest.Level) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.takeSemaphore()
		defer c.releazeSemaphore()
		arest.Debug("Pin: %d, Level: %s", pin, level.String())

		url := fmt.Sprintf("/digital/%d/%d\n\r", pin, level.Level())

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return err
		}

		var resp string
		select {
		case err := <-chErr:
			return err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		return nil
	}

	return err
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(ctx context.Context, pin int) (level arest.Level, err error) {
	select {
	case <-ctx.Done():
		return level, ctx.Err()
	default:
		c.takeSemaphore()
		defer c.releazeSemaphore()
		arest.Debug("Pin: %d", pin)

		url := fmt.Sprintf("/digital/%d\n\r", pin)
		data := make(map[string]interface{})

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return nil, err
		}

		var resp string
		select {
		case err := <-chErr:
			return nil, err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return nil, err
		}

		level = arest.NewLevel()
		if int(data["return_value"].(float64)) == arest.High {
			level.SetLevelHigh()
		} else {
			level.SetLevelLow()
		}

		return level, nil
	}

	return level, err

}

// ReadValue permit to read user variable
func (c *Client) ReadValue(ctx context.Context, name string) (value interface{}, err error) {
	select {
	case <-ctx.Done():
		return value, ctx.Err()
	default:
		c.takeSemaphore()
		defer c.releazeSemaphore()
		arest.Debug("Value name: %s", name)

		url := fmt.Sprintf("/%s\n\r", name)
		data := make(map[string]interface{})

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return nil, err
		}

		var resp string
		select {
		case err := <-chErr:
			return nil, err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return nil, err
		}

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
		c.takeSemaphore()
		defer c.releazeSemaphore()
		url := "/\n\r"
		data := make(map[string]interface{})

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return nil, err
		}

		var resp string
		select {
		case err := <-chErr:
			return nil, err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return nil, err
		}

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
		c.takeSemaphore()
		defer c.releazeSemaphore()
		arest.Debug("Function: %s, param: %s", name, param)

		url := fmt.Sprintf("/%s?params=%s\n\r", name, param)
		data := make(map[string]interface{})

		chErr := make(chan error)
		chRes := make(chan string)
		go c.read(ctx, chRes, chErr)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return value, err
		}

		var resp string
		select {
		case err := <-chErr:
			return value, err
		case resp = <-chRes:

		}

		arest.Debug("Resp: %s", resp)

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return value, err
		}

		if temp, ok := data["return_value"]; ok {
			value = int(temp.(float64))
		} else {
			err = errors.Errorf("Function %s not found", name)
		}

		return value, nil
	}

	return value, err

}

func (c *Client) read(ctx context.Context, chRes chan string, chErr chan error) {

	select {
	case <-ctx.Done():
		chErr <- ctx.Err()
	default:
		ch := make(chan bool)

		buffer := make([]byte, 2048)
		var resp strings.Builder

		ctx, _ = context.WithTimeout(ctx, c.timeout)

		go c.watchdog(ctx, ch, chErr)

		loop := true

		for loop {
			select {
			case <-chErr:
				return
			default:
				n, err := c.serialPort.Read(buffer)
				if err != nil {
					ch <- true
					chErr <- err
					return
				}
				if n == 0 {
					loop = false
					break
				}
				resp.Write(buffer[:n])

				if strings.Contains(string(buffer[:n]), "\n") {
					loop = false
					break
				}
			}

		}

		ch <- true
		chRes <- resp.String()
		return
	}

}

func (c *Client) takeSemaphore() {
	c.sem <- 1
}

func (c *Client) releazeSemaphore() {
	<-c.sem
}

func (c *Client) watchdog(ctx context.Context, ch chan bool, chErr chan error) {

	select {
	case <-ctx.Done():
		c.Client().Close()

		c.takeSemaphore()
		defer c.releazeSemaphore()

		serialPort, err := open(c.url)
		if err != nil {
			chErr <- err
			c.serialPort = nil
			return
		}
		c.serialPort = serialPort
		return
	case <-ch:
		return
	}
}

func open(url string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	serialPort, err := serial.Open(url, mode)
	if err != nil {
		originalErr := err
		ports, err := serial.GetPortsList()
		if err != nil {
			arest.Debug("%s", err)
		}
		if len(ports) == 0 {
			arest.Debug("No serial ports found!")
		}
		for _, port := range ports {
			arest.Debug("Found port: %v\n", port)
		}

		return nil, originalErr
	}

	time.Sleep(time.Second * 1)

	return serialPort, nil
}
