package svc

import (
	"fmt"
)

type NamedComponent struct {
	Component
	Name string
}

// NamedComponent manages a list of zero or more ordered components
// tagged with optional names. It is patterned off of the
// aws-sdk-go/aws/request/HandlerList type.
type NamedComponentList struct {
	list []*NamedComponent
}

func (l *NamedComponentList) Len() int {
	return len(l.list)
}

func (l *NamedComponentList) PushBack(c Component) {
	name := fmt.Sprintf("__anonymous%d", len(l.list))
	l.PushBackNamed(&NamedComponent{c, name})
}

func (l *NamedComponentList) PushBackNamed(c *NamedComponent) {
	if cap(l.list) == 0 {
		l.list = make([]*NamedComponent, 0, 5)
	}
	l.list = append(l.list, c)
}

func (l *NamedComponentList) PushFront(c *NamedComponent) {
	name := fmt.Sprintf("__anonymous%d", len(l.list))
	l.PushFrontNamed(&NamedComponent{c, name})
}

func (l *NamedComponentList) PushFrontNamed(c *NamedComponent) {
	if cap(l.list) == len(l.list) {
		// Allocating new List required
		l.list = append([]*NamedComponent{c}, l.list...)
	} else {
		// Enough room to prepend into list.
		l.list = append(l.list, &NamedComponent{})
		copy(l.list[1:], l.list)
		l.list[0] = c
	}
}

func (l *NamedComponentList) FindComponent(name string) Component {
	for _, named := range l.list {
		if named.Name == name {
			return named.Component
		}
	}
	return nil
}
