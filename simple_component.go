package svc

import (
	"context"

	"github.com/spf13/cobra"
)

// SimpleComponent is a concrete implementation of both Component and
// CommandInitializer allowing for simple components to be constructed
// inline out of handler functions.
type SimpleComponent struct {
	OnStart            func(ctx context.Context) error
	OnStop             func() error
	OnKill             func() error
	CommandInitializer CommandInitializer
}

func (s *SimpleComponent) Start(ctx context.Context) error {
	if s.OnStart != nil {
		return s.OnStart(ctx)
	}
	return nil
}

func (s *SimpleComponent) Stop() error {
	if s.OnStop != nil {
		return s.OnStop()
	}
	return nil
}

func (s *SimpleComponent) Kill() error {
	if s.OnKill != nil {
		return s.OnKill()
	}
	return nil
}

func (s *SimpleComponent) CommandInitialize(cmd *cobra.Command) {
	if s.CommandInitializer != nil {
		s.CommandInitialize(cmd)
	}
}
