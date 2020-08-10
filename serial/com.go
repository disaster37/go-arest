package serial

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ReadReady struct{}

func (c *Client) read() {
	buffer := make([]byte, 2048)
	var resp strings.Builder

	isReady := false
	for !isReady {
		_, err := c.serialPort.Write([]byte("/ready\n\r"))
		if err != nil {
			isReady = true
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	c.channel.Update(ReadReady{})
	log.Debugf("Reader on serial ready")
	for {
		n, err := c.serialPort.Read(buffer)
		log.Debugf("Receive: %v bytes", n)
		if err != nil {
			c.channel.Update(err)
		}
		if n == 0 {
			log.Debugf("Read: %s", resp.String())
			c.channel.Update(resp.String())
			resp.Reset()
		}
		resp.Write(buffer[:n])
		log.Debug(string(buffer[:n]))

		if strings.Contains(string(buffer[:n]), "\n") {
			log.Debugf("Read: %s", resp.String())
			c.channel.Update(resp.String())
			resp.Reset()
		}
	}

}

func (c *Client) takeSemaphore() {
	c.sem <- 1
}

func (c *Client) releazeSemaphore() {
	<-c.sem
}
