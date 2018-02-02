package cbp

import (
	"github.com/rs/xid"
)


type composite struct {
	id _id 
	components []*component 
	connections []*connection 
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

func addConnection(c *composite, upstreamComponentUID string, downstreamComponentUID string) error {
	var (
		u *component
		d *component 
	)
	for _, v := range c.components {
		if upstreamComponentUID == v.id.uid {
			u = v 
		} else if downstreamComponentUID == v.id.uid {
			d = v 
		}
	}
	cn, err := newConnection(u.id.name+"/"+d.id.name, u, d)
	if err != nil {
		return err 
	}
	cn.id.uid = xid.New().String()
	c.connections = append(c.connections, cn)
	return nil 
}

// AddInterConnection is for connections to outside of this c
// ..(c to component or c to c)
// func AddInterConnection() {}

// RunComposite blah
func RunComposite(c *composite) {
	// todo: validate c
	c.head = c.connections[0].upstreamComponent
	c.tail = c.connections[len(c.connections)-1].downstreamComponent
	for _, cn := range c.connections {
		runComponent(cn.upstreamComponent)
		runComponent(cn.downstreamComponent)
	}
}