package cbp

import (
	"github.com/rs/xid"
)

// Composite blah
type Composite struct {
	id _id 
	components []*Component 
	connections []*connection 
	head *Component 
	tail *Component 
}

// NewComposite retuerns a new, empty Composite object
func NewComposite(name string) *Composite {
	c := new(Composite)
	c.id.name = name 
	c.id.uid = xid.New().String()
	return c
}

// AddComponent blah
func (c *Composite) AddComponent(comp *Component) {
	c.components = append(c.components, comp)
}

// AddConnection blah
func (c *Composite) AddConnection(upstreamComponentUID string, downstreamComponentUID string) error {
	var (
		u *Component
		d *Component 
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
// ..(c to Component or c to c)
// func AddInterConnection() {}

// RunComposite blah
func RunComposite(c *Composite) {
	// todo: validate c
	c.head = c.connections[0].upstreamComponent
	c.tail = c.connections[len(c.connections)-1].downstreamComponent
	for _, cn := range c.connections {
		cn.upstreamComponent.RunComponent()
		cn.downstreamComponent.RunComponent()
	}
}