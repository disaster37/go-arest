package serial

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

type ReadReady struct{}

func read(c *Client) {
	buffer := make([]byte, 2048)
	var resp strings.Builder

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
		}
		resp.Write(buffer[:n])
		log.Debug(string(buffer[:n]))

		if strings.Contains(string(buffer[:n]), "\n") {
			log.Debugf("Read: %s", resp.String())
			c.channel.Update(resp.String())
		}
	}

}

func (c *Client) takeSemaphore() {
	c.sem <- 1
}

func (c *Client) releazeSemaphore() {
	<-c.sem
}
