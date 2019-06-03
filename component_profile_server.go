package svc

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

type ProfileServer struct {
	addr   string
	enable bool
	server *manners.GracefulServer
}

func (p *ProfileServer) CommandInitialize(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&p.enable, "profile-server-enable", "", false,
		"If enabled, start an HTTP profile server for diagnostics")
	cmd.Flags().StringVarP(&p.addr, "profile-server-addr", "", ":8081",
		"Address to bind the HTTP profile server to, if enabled")
}

func (p *ProfileServer) Start(ctx context.Context) error {
	if p.enable {
		log.WithFields(log.Fields{
			"action": "profile_server",
			"status": "start",
			"addr":   p.addr,
		}).Info()

		exp.Exp(metrics.DefaultRegistry)
		metricscharts.Register()

		p.server = manners.NewWithServer(&http.Server{
			Addr:    p.addr,
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
	return nil
}

func (p *ProfileServer) Kill() error {
	return p.Stop()
}
