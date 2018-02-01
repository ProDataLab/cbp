/*
	TODO:
		- does "new()" return an error signal as well?
*/

package cbp

import "github.com/rs/xid"

type Component struct {
	id            Id
	inSockets     []*Socket
	outSockets    []*Socket
	configSockets []*Socket
	reportSockets []*Socket
}

// NewComonent allocates and returns a new Component to the user.
func NewComponent(name string) (*Component, error) {
	c := new(Component)
	c.id.name = name
	c.id.uid = xid.New().String()
	return c, nil 
}

func AddSocket(name string, socketType string, transportType string, url string) {
	socket.newSocket(name, socketType, transportType, url)
}
