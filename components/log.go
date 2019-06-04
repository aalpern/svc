package components

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LogConfigComponent struct {
	verbose bool
}

func (p *LogConfigComponent) CommandInitialize(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&p.verbose, "log-verbose", "v", false,
		"Set logging level to verbose (debug)")
}

func (p *LogConfigComponent) Start(ctx context.Context) error {
	if p.verbose {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func (p *LogConfigComponent) Stop() error {
	return nil
}

func (p *LogConfigComponent) Kill() error {
	return nil
}
