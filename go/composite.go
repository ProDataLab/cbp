package cbp

import (
	"github.com/rs/xid"
)


type Composite struct {
	id Id 
	components []*Component 
	connections []*Connection 
	head *Component 
	tail *Component 
}

// NewComposite retuerns a new, empty Composite object
func NewComposite(name string) *Composite, error {
	c := new(Composite)
	c.id.name = name 
	c.id.uid = xid.New().String()
	return c, nil 
}

func AddComponent(composite *Composite, component *Component) {
	composite.components = append(composite.components, component)
}

func AddIntraConnection(composite *Composite, upstreamComponentUid string, downstreamComponentUid string) {
	var (
		u *Component
		d *Component 
	)
	for _, v := range composite.components {
		if upstreamComponentUid == v.id.uid {
			u = v 
		} else if downstreamComponentUid == v.id.uid {
			d = v 
		}
	}
	c := NewConnection(u.id.name+"/"+d.id.name, u, d)
	d.id.uid = xid.New().String()
	composite.connections = append(composite.connections, c)
}

// AddInterConnection is for connections to outside of this composite
// ..(composite to component or composite to composite)
// func AddInterConnection() {}

func RunComposite(composite *Composite) error {
	// todo: validate composite
	composite.head = composite.connections[0]
	composite.tail = composite.connections[len(composite.connections)-1]
	for i, c := range composite.connections {
		RunComponent(c.upstreamComponent)
		RunComponent(c.downstreamComponent)
	}
}