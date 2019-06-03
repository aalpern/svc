package svc

import (
	"github.com/spf13/cobra"
)

type CommandInitializer interface {
	CommandInitialize(cmd *cobra.Command)
}
