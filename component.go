/*
	TODO:
		- does "new()" return an error signal as well?
		- think about reporting sockets and best way to
			implement these
		- ..and what about config(sockets) ?
*/

package cbp

import (
	"fmt"

	"github.com/rs/xid"
)

// Component provides the structure of a Component
type Component struct {
	id            _id
	inSockets     []*socket
	outSockets    []*socket
	configSockets []*socket
	reportSockets []*socket
	inChannel     chan []byte
	outChannel    chan []byte
}

// NewComponent allocates and returns a new Component to the user.
func NewComponent(name string) (*Component, error) {
	c := new(Component)
	c.id.name = name
	c.id.uid = xid.New().String()
	return c, nil
}

// AddSocket adds a socket to the component
func (c *Component) AddSocket(name string, st SocketType, tt TransportType, url string) error {
	// TODO: check both SocketType and transportType
	s, err := newSocket(name, st, tt, url)
	if err != nil {
		return err
	}
	switch st {
	case "req":
		c.outSockets = append(c.outSockets, s)
	case "rep":
		c.inSockets = append(c.inSockets, s)
	case "pub":
		c.outSockets = append(c.outSockets, s)
	case "sub":
		c.inSockets = append(c.inSockets, s)
	case "push":
		c.outSockets = append(c.outSockets, s)
	case "pull":
		c.inSockets = append(c.inSockets, s)
	}
	return nil
}

// RunComponent blah
func (c *Component) RunComponent() error {
	var all [][]*socket
	all = append(all, c.inSockets)
	all = append(all, c.outSockets)
	all = append(all, c.configSockets)
	all = append(all, c.reportSockets)
	for _, r := range all {
		for _, s := range r {
			s.run()
		}
	}
	c.inChannel = make(chan []byte)
	c.outChannel = make(chan []byte)
	fanIn := func(ch chan []byte) {
		for msg := range ch {
			c.inChannel <- msg
		}
	}
	fanOut := func(sc chan []byte) {
		for msg := range c.outChannel {
			fmt.Println(string(msg))
			sc <- msg
		}
	}
	for _, s := range c.inSockets {
		fmt.Println(s.id.name)
		go fanIn(s.recvChannel)
	}
	for _, s := range c.outSockets {
		fmt.Println(s.id.name)
		go fanOut(s.sendChannel)
	}
	return nil
}

// InChannel blah
func (c *Component) InChannel() chan []byte {
	return c.inChannel
}

// OutChannel blah
func (c *Component) OutChannel() chan []byte {
	return c.outChannel
}
