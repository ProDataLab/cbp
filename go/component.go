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

type component struct {
	id            Id
	inSockets     []*socket
	outSockets    []*socket
	configSockets []*socket
	reportSockets []*socket
	inChannel 		chan []byte
	outChannel 		chan []byte
}

// newComonent allocates and returns a new component to the user.
func newComponent(name string) (*component, error) {
	c := new(component)
	c.id.name = name
	c.id.uid = xid.New().String()
	return c, nil
}

func addSocket(name string, comp *component, socketType string, transportType string, url string) {
	// TODO: check both socketType and transportType
	if s, err := socket.newSocket(name, socketType, transportType, url); err != nil {
		return err
	}
	c := comp
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

func runComponent(comp *component) error {

}