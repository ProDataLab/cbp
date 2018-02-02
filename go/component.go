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
	id            _id
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

func addSocket(name string, comp *component, st socketType, tt transportType, url string) error {
	// TODO: check both socketType and transportType
	s, err := newSocket(name, st, tt, url)
	if err != nil {
		return err
	}
	c := comp
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

func runComponent(comp *component) error {
	return nil
}