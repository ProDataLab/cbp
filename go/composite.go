package cbp

import (
	"github.com/rs/xid"
)


type composite struct {
	id Id 
	components []*component 
	connections []*Connection 
	head *component 
	tail *component 
}

// newComposite retuerns a new, empty composite object
func newComposite(name string) *composite {
	c := new(composite)
	c.id.name = name 
	c.id.uid = xid.New().String()
	return c
}

func addComponent(c *composite, comp *component) {
	c.components = append(c.components, comp)
}

func addConnection(c *composite, upstreamComponentUid string, downstreamComponentUid string) {
	var (
		u *component
		d *component 
	)
	for _, v := range c.components {
		if upstreamComponentUid == v.id.uid {
			u = v 
		} else if downstreamComponentUid == v.id.uid {
			d = v 
		}
	}
	cn := NewConnection(u.id.name+"/"+d.id.name, u, d)
	cn.id.uid = xid.New().String()
	c.connections = append(c.connections, cn)
}

// AddInterConnection is for connections to outside of this c
// ..(c to component or c to c)
// func AddInterConnection() {}

func RunComposite(c *composite) error {
	// todo: validate c
	c.head = c.connections[0]
	c.tail = c.connections[len(c.connections)-1]
	for i, cn := range c.connections {
		RunComponent(cn.upstreamComponent)
		RunComponent(cn.downstreamComponent)
	}
}