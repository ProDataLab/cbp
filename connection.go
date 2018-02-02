package cbp

import (
	"github.com/rs/xid"
)


type connection struct {
	id _id
	upstreamComponent *component 
	downstreamComponent *component 
}

// newConnection returns a new Connection object
func newConnection(name string, upstreamComponent *component, downstreamComponent *component) (*connection, error) {
	// todo: type check component
	conn := new(connection)
	conn.id.name = name 
	conn.id.uid = xid.New().String()
	conn.upstreamComponent = upstreamComponent
	conn.downstreamComponent = downstreamComponent
	return conn, nil 
}