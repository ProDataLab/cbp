package cbp

import (
	"github.com/rs/xid"
)


type Connection struct {
	id Id
	upstreamComponent *Component 
	downstreamComponent *Component 
}

// NewConnection returns a new Connection object
func NewConnection(name string, upstreamComponent *Component, downstreamComponent *Component) (*Connection, error) {
	// todo: type check component
	c := new(Connection)
	c.id.name = name 
	c.id.uid = xid.New().String()
	c.upstreamComponent = upstreamComponent
	c.downstreamComponent = downstreamComponent
	return c, nil 
}