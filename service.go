package svc

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Service defines the root command and entry point for a persistent
// service process. The service instance controls overall startup and
// shutdown flow, and with provided default configuration options
// provides standardized configuration and setup of common
// requirements for well-behavced services, such as logging and
// responding to signals.
type Service struct {
	*cobra.Command

	// Name is an identifier for the service which will be used in
	// command line help, logging, and metrics.
	Name string

	// Global is a component which will be initialized for every
	// command implemented by the service - common initialization such
	// as logging is typically handled here.
	Global Component

	exit chan int
}

type contextKey int

const serviceKey contextKey = iota

// WithServiceContext returns a derived context with the supplied
// pointer to a Service instance bound as a value.
func WithServiceContext(ctx context.Context, svc *Service) context.Context {
	return context.WithValue(ctx, serviceKey, svc)
}

// GetService extracts and returns the service pointer bound as a
// context value, or nil if no service pointer was found in the
// supplied context.
func GetService(ctx context.Context) *Service {
	if svc, ok := ctx.Value(serviceKey).(*Service); ok {
		return svc
	}
	return nil
}

type ServiceConfig func(svc *Service) error

// NewService constructs a new Service instance from the given name,
// descript, and configs.
func NewService(name, description string, configs ...ServiceConfig) (*Service, error) {
	svc := &Service{
		Command: NewCommand(name, description),
		Name:    name,
		exit:    make(chan int),
	}

	// Persistent pre run handler initializes the global component for
	// all commands
	svc.Command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ctx := WithServiceContext(context.Background(), svc)
		if err := svc.Start(ctx); err != nil {
			log.WithFields(log.Fields{
				"action": "service_global_start",
				"status": "error",
				"error":  err,
			}).Error("Error starting global component")
			return err
		}
		return nil
	}

	// And a persistent post run handler stops the global component
	// after all commands
	svc.Command.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		if err := svc.Stop(); err != nil {
			log.WithFields(log.Fields{
				"action": "service_global_stop",
				"status": "error",
				"error":  err,
			}).Error("Error stopping global component")
			return err
		}
		return nil
	}

	for _, cfg := range configs {
		if err := cfg(svc); err != nil {
			return nil, err
		}
	}

	if ci, ok := svc.Global.(CommandInitializer); ok {
		ci.CommandInitialize(svc.Command)
	}

	return svc, nil
}

func (svc *Service) Start(ctx context.Context) error {
	if svc.Global != nil {
		log.WithFields(log.Fields{
			"action": "service_start",
			"status": "starting_global",
		}).Debug()
		return svc.Global.Start(ctx)
	}
	return nil
}

func (svc *Service) Stop() error {
	if svc.Global != nil {
		log.WithFields(log.Fields{
			"action": "service_stop",
			"status": "stopping_global",
		}).Debug()
		return svc.Global.Stop()
	}
	return nil
}

func (svc *Service) FindComponent(name string) Component {
	if svc.Global != nil {
		if cc, ok := svc.Global.(*CompositeComponent); ok {
			return cc.FindComponent(name)
		}
	}
	return nil
}

func (svc *Service) Exit(code int) {
	svc.exit <- code
}

func WithCommand(cmd ...*cobra.Command) ServiceConfig {
	return func(svc *Service) error {
		svc.Command.AddCommand(cmd...)
		return nil
	}
}

func WithLongDescription(desc string) ServiceConfig {
	return func(svc *Service) error {
		svc.Command.Long = desc
		return nil
	}
}

func WithGlobal(cmp Component) ServiceConfig {
	return func(svc *Service) error {
		svc.Global = cmp
		return nil
	}
}

// WithCommandHandler creates a new command for the service process
// and binds it to the supplied Component, whose Start() method
// becomes the main loop of the command.
//
// The Start() method of the handler component will be run in a new go
// routine, while the main go routine waits for an exit code on a
// shutdown channel. The Start() routine of the main handler component
// may block (e.g. if it implements a network server).
//
// A PostRun handler on the command handles calling Stop() to shut
// down the component, with a configurable timeout before calling
// Kill() if the ordered shutdown takes too long.
func WithCommandHandler(name, description string, handler Component) ServiceConfig {
	return func(svc *Service) error {
		killTimeout := 30 * time.Second
		cmd := NewCommand(name, description, handler)

		cmd.Flags().DurationVar(&killTimeout, "service-kill-timeout", 30*time.Second,
			"Time to wait for ordered shutdown to complete before hard exit")

		cmd.Run = func(cmd *cobra.Command, args []string) {
			log.WithFields(log.Fields{
				"action": "command_handler",
				"status": "start",
			}).Info()
			go func() {
				ctx := WithServiceContext(context.Background(), svc)
				if err := handler.Start(ctx); err != nil {
					log.WithFields(log.Fields{
						"action": "command_handler",
						"status": "start_error",
						"error":  err,
					}).Error("Error starting command handler")
					svc.Exit(-1)
				}
			}()
			code := <-svc.exit
			log.WithFields(log.Fields{
				"action":    "command_handler",
				"status":    "done",
				"exit_code": code,
			}).Info()
		}

		// PostRun handler issues the call to Stop() to perform
		// ordered shut down, using some go routine/channel
		// shenangigans to implement a timeout.
		cmd.PostRun = func(cmd *cobra.Command, args []string) {
			stop := make(chan interface{})
			go func() {
				if err := handler.Stop(); err != nil {
					log.WithFields(log.Fields{
						"action": "command_handler",
						"status": "stop_error",
						"error":  err,
					}).Error("Error stopping command handler")
				}
				// We don't send any data back on the channel, just
				// close it to indicate that the stop has completed.
				close(stop)
			}()

			select {
			case <-stop:
			case <-time.After(killTimeout):
				log.WithFields(log.Fields{
					"action":  "command_handler",
					"status":  "stop_timeout",
					"timeout": killTimeout,
				}).Info("Ordered stop timed out, killing")
				handler.Kill()
			}
		}

		svc.Command.AddCommand(cmd)

		return nil
	}
}
