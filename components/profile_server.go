package components

import (
	"context"
	"net/http"

	"github.com/aalpern/go-metrics-charts"
	"github.com/braintree/manners"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	DefaultProfileAddr = ":8081"
)

type ProfileServer struct {
	Addr   string
	Enable bool
	server *manners.GracefulServer
}

func (p *ProfileServer) CommandInitialize(cmd *cobra.Command) {
	defaultaddr := p.Addr
	if defaultaddr == "" {
		defaultaddr = DefaultProfileAddr
	}
	cmd.Flags().BoolVar(&p.Enable, "profile-server-enable", p.Enable,
		"If enabled, start an HTTP profile server for diagnostics")
	cmd.Flags().StringVar(&p.Addr, "profile-server-addr", defaultaddr,
		"Address to bind the HTTP profile server to, if enabled")
}

func (p *ProfileServer) Start(ctx context.Context) error {
	if p.Enable {
		log.WithFields(log.Fields{
			"action": "profile_server",
			"status": "start",
			"addr":   p.Addr,
		}).Info()

		exp.Exp(metrics.DefaultRegistry)
		metricscharts.Register()

		p.server = manners.NewWithServer(&http.Server{
			Addr:    p.Addr,
			Handler: http.DefaultServeMux,
		})

		go func() {
			if err := p.server.ListenAndServe(); err != nil {
				log.WithFields(log.Fields{
					"action": "profile_server",
					"status": "error",
					"error":  err,
				}).Error("Profile server exited with error")
			} else {
				log.WithFields(log.Fields{
					"action": "profile_server",
					"status": "done",
				}).Info()
			}
		}()
	}

	return nil
}

func (p *ProfileServer) Stop() error {
	if p.server != nil {
		p.server.BlockingClose()
	}
	return nil
}

func (p *ProfileServer) Kill() error {
	if p.server != nil {
		p.server.Close()
	}
	return nil
}
