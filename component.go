package svc

import (
	"context"
	"fmt"
)

// Component defines the basic lifecycle interface for types that can
// be started up and shut down.
type Component interface {
	Start(ctx context.Context) error
	Stop() error
	Kill() error
}

type NamedComponent struct {
	Component
	Name string
}

// NamedComponent manages a list of zero or more ordered components
// tagged with optional names. It is patterned off of the
// aws-sdk-go/aws/request/HandlerList type.
type NamedComponentList struct {
	list []NamedComponent
}

func (l *NamedComponentList) Len() int {
	return len(l.list)
}

func (l *NamedComponentList) PushBack(c Component) {
	name := fmt.Sprintf("__anonymous%d", len(l.list))
	l.PushBackNamed(NamedComponent{c, name})
}

func (l *NamedComponentList) PushBackNamed(c NamedComponent) {
	if cap(l.list) == 0 {
		l.list = make([]NamedComponent, 0, 5)
	}
	l.list = append(l.list, c)
}

func (l *NamedComponentList) PushFront(c NamedComponent) {
	name := fmt.Sprintf("__anonymous%d", len(l.list))
	l.PushFrontNamed(NamedComponent{c, name})
}

func (l *NamedComponentList) PushFrontNamed(c NamedComponent) {
	if cap(l.list) == len(l.list) {
		// Allocating new List required
		l.list = append([]NamedComponent{c}, l.list...)
	} else {
		// Enough room to prepend into list.
		l.list = append(l.list, NamedComponent{})
		copy(l.list[1:], l.list)
		l.list[0] = c
	}
}

func (l *NamedComponentList) FindComponent(name string) Component {
	for _, named := range l.list {
		if named.Name == Name {
			return named.Component
		}
	}
	return nil
}

// CompositeComponent is an implementation of Component for composing
// multiple components together. Implementations of all the service
// framework interfaces that components can possess will be forwarded
// to the children if they implement those interfaces.
type CompositeComponent struct {
	children NamedComponentList
}

type CompositeComponentOption func(*CompositeComponent) error

func NewCompositeComponent(opts ...CompositeComponentOption) (*CompositeComponent, error) {
	c := &CompositeComponent{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func WithComponent(c Component) CompositeComponentOption {
	return func(cc *CompositeComponent) error {
		cc.Components.PushBack(c)
		return nil
	}
}

func WithNamedComponent(name string, c Component) CompositeComponentOption {
	return func(cc *CompositeComponent) error {
		cc.Components.PushBackNamed(name, c)
		return nil
	}
}

func (c *CompositeComponent) Start(ctx context.Context) error {
	for _, child := range c.children.list {
		if err := child.Component.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeComponent) Stop() error {
	for _, child := range c.children.list {
		if err := child.Component.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeComponent) Kill() error {
	for _, child := range c.children.list {
		if err := child.Component.Kill(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeComponent) FindComponent(name string) Component {
	return c.children.FindComponent(name)
}

func (c *CompositeComponent) CommandInitialize(cmd *cobra.Command) {
	for _, child := range s.children.list {
		if ci, ok := child.Component.(CommandInitializer); ok {
			ci.CommandInitialize(cmd)
		}
	}
}
