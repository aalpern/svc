package svc

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	ShutdownSignals = []os.Signal{
		syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
	}
)

type SignalHandler func(context.Context, os.Signal)

func NewSignalWatcher(handler SignalHandler, signals ...os.Signal) Component {
	return &SimpleComponent{
		OnStart: func(ctx context.Context) error {
			go func() {
				chn := make(chan os.Signal, 1)
				signal.Notify(chn, signals...)
				for {
					sig := <-chn
					log.WithFields(log.Fields{
						"action": "signal",
						"signal": sig,
					}).Debug("Trapped signal")
					handler(ctx, sig)
				}
			}()
			return nil
		},
	}
}

func WithShutdownWatcher() CompositeComponentOption {
	return WithNamedComponent("shutdown-watcher",
		NewSignalWatcher(func(ctx context.Context, sig os.Signal) {
			log.WithFields(log.Fields{
				"action": "shutdown_signal",
				"status": "signaled",
				"signal": sig,
			}).Info("Initiating shutdown")
			if svc := GetService(ctx); svc != nil {
				svc.Exit(0)
			} else {
				log.WithFields(log.Fields{
					"action": "shutdown_signal",
					"status": "no_service",
				}).Warn("No bound Service instance, performing hard exit")
				os.Exit(0)
			}
		}, ShutdownSignals...))
}
