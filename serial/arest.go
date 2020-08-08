package rest

import (
	"encoding/json"
	"fmt"

	"github.com/disaster37/go-arest"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// Client implement arest interface
type Client struct {
	serialPort serial.Port
}

// NewClient permit to initialize new client Object
func NewClient(url string) (arest.Arest, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	serialPort, err := serial.Open(url, mode)
	if err != nil {
		log.Errorf("Error appear when open serial port, we scrut serial port")
		originalErr := err
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Error(err)
		}
		if len(ports) == 0 {
			log.Info("No serial ports found!")
		}
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
		}

		return nil, originalErr
	}

	return &Client{
		serialPort: serialPort,
	}, nil
}

// Client permit to get curent resty client
func (c *Client) Client() serial.Port {
	return c.serialPort
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(pin int, mode arest.Mode) (err error) {

	log.Debugf("Pin: %d, Mode: %s", pin, mode.String())

	url := fmt.Sprintf("/mode/%d/%s\n\r", pin, mode.Mode())

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return err
	}

	log.Debugf("Sent: %v bytes", resp)

	return nil

}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(pin int, level arest.Level) (err error) {

	log.Debugf("Pin: %d, Level: %s", pin, level.String())

	url := fmt.Sprintf("/digital/%d/%d\n\r", pin, level.Level())

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return err
	}

	log.Debugf("Sent: %v bytes", resp)

	return err
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(pin int) (level arest.Level, err error) {

	log.Debugf("Pin: %d", pin)

	url := fmt.Sprintf("/digital/%d\n\r", pin)
	data := make(map[string]interface{})

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}
	log.Debugf("Sent: %v bytes", resp)

	buffer := make([]byte, 0, 0)
	resp, err = c.serialPort.Read(buffer)
	if err != nil {
		return nil, err
	}
	log.Debugf("Receive: %v bytes", resp)

	err = json.Unmarshal(buffer, &data)
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

	log.Debugf("Value name: %s", name)

	url := fmt.Sprintf("/%s\n\r", name)
	data := make(map[string]interface{})
	buffer := make([]byte, 0, 0)

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}
	log.Debugf("Sent: %v bytes", resp)

	resp, err = c.serialPort.Read(buffer)
	if err != nil {
		return nil, err
	}
	log.Debugf("Receive: %v bytes", resp)

	err = json.Unmarshal(buffer, &data)
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

	url := "/\n\r"
	data := make(map[string]interface{})
	buffer := make([]byte, 0, 0)

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return nil, err
	}
	log.Debugf("Sent: %v bytes", resp)

	resp, err = c.serialPort.Read(buffer)
	if err != nil {
		return nil, err
	}
	log.Debugf("Receive: %v bytes", resp)

	err = json.Unmarshal(buffer, &data)
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

	log.Debugf("Function: %s, param: %s", name, param)

	url := fmt.Sprintf("/%s?params=%s\n\r", name, param)
	data := make(map[string]interface{})
	buffer := make([]byte, 0, 0)

	resp, err := c.serialPort.Write([]byte(url))
	if err != nil {
		return value, err
	}
	log.Debugf("Sent: %v bytes", resp)

	resp, err = c.serialPort.Read(buffer)
	if err != nil {
		return value, err
	}
	log.Debugf("Receive: %v bytes", resp)

	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return value, err
	}

	if temp, ok := data["return_value"]; ok {
		value = int(temp.(float64))
	} else {
		errors.Errorf("Function %s not found", name)
	}

	return value, err
}
