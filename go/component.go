/*
	TODO:
		- does "new()" return an error signal as well?
		- think about reporting sockets and best way to
			implement these
		- ..and what about config(sockets) ?
*/

package cbp

import (
	"github.com/rs/xid"
)

type Component struct {
	id            Id
	inSockets     []*Socket
	outSockets    []*Socket
	configSockets []*Socket
	reportSockets []*Socket
	inChannel 		[]byte chan 
	outChannel 		[]byte chan 
}

// NewComonent allocates and returns a new Component to the user.
func NewComponent(name string) (*Component, error) {
	c := new(Component)
	c.id.name = name
	c.id.uid = xid.New().String()S
	return c, nil
}

func AddSocket(name string, component *Component, socketType string, transportType string, url string) {
	// TODO: check both socketType and transportType
	if s, err := socket.newSocket(name, socketType, transportType, url); err != nil {
		return err
	}
	c := component
	switch socketType {
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
	return s, nil
}

func Run(component *Component) error {

}