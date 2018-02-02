package cbp

import (
	"github.com/rs/xid"
)


type connection struct {
	id _id
	upstreamComponent *Component 
	downstreamComponent *Component 
}

// newConnection returns a new Connection object
func newConnection(name string, upstreamComponent *Component, downstreamComponent *Component) (*connection, error) {
	// todo: type check Component
	conn := new(connection)
	conn.id.name = name 
	conn.id.uid = xid.New().String()
	conn.upstreamComponent = upstreamComponent
	conn.downstreamComponent = downstreamComponent
	return conn, nil 
}