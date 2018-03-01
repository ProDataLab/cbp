/*
	TODO:
		- does "new()" return an error signal as well?
		- think about reporting sockets and best way to
			implement these
		- ..and what about config(sockets) ?
*/

package cbp

import (
	"errors"
	"fmt"
	"strings"

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

var (
	// ErrComponentHasNoSockets is returned when the component is run with no sockets attached.
	ErrComponentHasNoSockets = errors.New("this component has no sockets")
)

// NewComponent allocates and returns a new Component to the user.
func NewComponent(name string) (*Component, error) {
	c := new(Component)
	c.id.name = name
	c.id.uid = xid.New().String()
	return c, nil
}

// Name returns the name of this component
func (c *Component) Name() string {
	return c.id.name
}

// AddSocket adds a socket to the component
func (c *Component) AddSocket(name string, st SocketType, tt TransportType, url string) error {
	// TODO: check both SocketType and transportType
	fmt.Printf("INFO: In c.AddSocket.. socketType: %s\n", string(st))
	s, err := newSocket(name, st, tt, url)
	if err != nil {
		return err
	}
	if name == "config" {
		c.configSockets = append(c.configSockets, s)
		err = s.setSubscriptionFilters(c.Name+"-config")
		if err != nil {
			return err 
		}
		return nil 
	}
	if strings.Contains(name, "-report") {
		c.reportSockets = append(c.reportSockets, s)
		return nil 
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

// AddConfigSocket is used to configure the component. A component is initially started with
// only this socket. Further configuration is done dynamically.
func (c *Component) AddConfigSocket(url string) error {
	return c.AddSocket("config", "sub", "tcp", url)
}

// AddReportSocket is used to report all errors, etc from the component. It is a pub socket
// so any sink should subscribe to it.
func (c *Component) AddReportSocket(reportComponentName string, url string) error {
	return c.AddSocket(reportComponentName+"-report", "push", "tcp", url)
}

// RunComponent blah
func (c *Component) Run() error {
	var all [][]*socket
	all = append(all, c.configSockets)
	all = append(all, c.reportSockets)
	if c.inSockets != nil {
		all = append(all, c.inSockets)
	}
	if c.outSockets != nil {
		all = append(all, c.outSockets)
	}
	if len(all) == 0 {
		return ErrComponentHasNoSockets
	}
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
		switch st {
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

// // InChannel blah
// func (c *Component) InChannel() chan []byte {
// 	return c.inChannel
// }

// // OutChannel blah
// func (c *Component) OutChannel() chan []byte {
// 	return c.outChannel
// }

// Send sends the msgpack encoded byte array to downstream
func (c *Component) Send(val []byte) {
	c.outChannel <- val
}

// Recv receives the msgpack encoded byte array from upstream
func (c *Component) Recv() []byte {
	return <-c.inChannel
}
