package serial

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/disaster37/go-arest"
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
func (c *Client) SetPinMode(pin int, mode arest.Mode) (err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	arest.Debug("Pin: %d, Mode: %s", pin, mode.String())

	url := fmt.Sprintf("/mode/%d/%s\n\r", pin, mode.Mode())

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return err
	}

	resp, err := c.read()
	if err != nil {
		return err
	}

	arest.Debug("Resp: %s", resp)

	return nil

}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(pin int, level arest.Level) (err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	arest.Debug("Pin: %d, Level: %s", pin, level.String())

	url := fmt.Sprintf("/digital/%d/%d\n\r", pin, level.Level())

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return err
	}

	resp, err := c.read()
	if err != nil {
		return err
	}

	arest.Debug("Resp: %s", resp)

	return err
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(pin int) (level arest.Level, err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	arest.Debug("Pin: %d", pin)

	url := fmt.Sprintf("/digital/%d\n\r", pin)
	data := make(map[string]interface{})

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}

	resp, err := c.read()
	if err != nil {
		return nil, err
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

	return level, err
}

// ReadValue permit to read user variable
func (c *Client) ReadValue(name string) (value interface{}, err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	arest.Debug("Value name: %s", name)

	url := fmt.Sprintf("/%s\n\r", name)
	data := make(map[string]interface{})

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}

	resp, err := c.read()
	if err != nil {
		return nil, err
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

// ReadValues permit to read user variable
func (c *Client) ReadValues() (values map[string]interface{}, err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	url := "/\n\r"
	data := make(map[string]interface{})

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}

	resp, err := c.read()
	if err != nil {
		return nil, err
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

// CallFunction permit to call user function
func (c *Client) CallFunction(name string, param string) (value int, err error) {
	c.takeSemaphore()
	defer c.releazeSemaphore()
	arest.Debug("Function: %s, param: %s", name, param)

	url := fmt.Sprintf("/%s?params=%s\n\r", name, param)
	data := make(map[string]interface{})

	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return value, err
	}

	resp, err := c.read()
	if err != nil {
		return value, err
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

	return value, err
}

func (c *Client) read() (string, error) {
	buffer := make([]byte, 2048)
	var resp strings.Builder
	finished := false

	go c.watchdog(&finished)

	for {
		n, err := c.serialPort.Read(buffer)
		if err != nil {
			finished = true
			return "", err
		}
		if n == 0 {
			break
		}
		resp.Write(buffer[:n])

		if strings.Contains(string(buffer[:n]), "\n") {
			break
		}
	}

	finished = true

	return resp.String(), nil
}

func (c *Client) takeSemaphore() {
	c.sem <- 1
}

func (c *Client) releazeSemaphore() {
	<-c.sem
}

func (c *Client) watchdog(finished *bool) {
	time.Sleep(c.timeout)
	if !*finished {
		c.Client().Close()
		// Try to reopen it
		go func() {

			isConnected := false
			c.takeSemaphore()
			defer c.releazeSemaphore()

			for !isConnected {
				serialPort, err := open(c.url)
				if err != nil {
					arest.Debug("Error when try to reconnect on serial port: %s", err.Error())
					time.Sleep(1 * time.Second)
				} else {
					arest.Debug("Successfully reopened serial port")
					c.serialPort = serialPort
					isConnected = true
				}
			}
		}()
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
