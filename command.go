package svc

import (
	"github.com/spf13/cobra"
)

type CommandInitializer interface {
	CommandInitialize(cmd *cobra.Command)
}

type CommandInitializerFn func(cmd *cobra.Command)

func (cif CommandInitializerFn) CommandInitialize(cmd *cobra.Command) {
	cif(cmd)
}

func NewCommand(name, description string, initializers ...interface{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: description,
	}
	for _, i := range initializers {
		if ci, ok := i.(CommandInitializer); ok {
			ci.CommandInitialize(cmd)
		}
	}
	return cmd
}
