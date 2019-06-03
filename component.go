package svc

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Component defines the basic lifecycle interface for types that can
// be started up and shut down.
type Component interface {
	Start(ctx context.Context) error
	Stop() error
	Kill() error
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
		cc.children.PushBack(c)
		return nil
	}
}

func WithNamedComponent(name string, c Component) CompositeComponentOption {
	return func(cc *CompositeComponent) error {
		cc.children.PushBackNamed(&NamedComponent{c, name})
		return nil
	}
}

func (c *CompositeComponent) Start(ctx context.Context) error {
	for i, child := range c.children.list {
		log.WithFields(log.Fields{
			"action":          "start_composite_component",
			"component_index": i,
			"component_name":  child.Name,
		}).Debug()
		if err := child.Component.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeComponent) Stop() error {
	for i := len(c.children.list) - 1; i >= 0; i-- {
		child := c.children.list[i]
		log.WithFields(log.Fields{
			"action":          "stop_composite_component",
			"component_index": i,
			"component_name":  child.Name,
		}).Debug()
		if err := child.Component.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompositeComponent) Kill() error {
	for i := len(c.children.list) - 1; i >= 0; i-- {
		child := c.children.list[i]
		log.WithFields(log.Fields{
			"action":          "kill_composite_component",
			"component_index": i,
			"component_name":  child.Name,
		}).Debug()
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
	for _, child := range c.children.list {
		if ci, ok := child.Component.(CommandInitializer); ok {
			ci.CommandInitialize(cmd)
		}
	}
}
