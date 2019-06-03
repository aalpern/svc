package httpsvc

import (
	"context"
	"net/http"

	"github.com/aalpern/svc"
	"github.com/spf13/cobra"
)

type HttpServiceComponent struct {
	addr     string
	handlers []http.Handler
}

type Option func(c *HttpServiceComponent)

func New(opts ...Option) *HttpServiceComponent {
	c := &HttpServiceComponent{
		addr:     ":8080",
		handlers: make([]http.Handler, 0, 8),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithHandler(h http.Handler) Option {
	return func(c *HttpServiceComponent) {
		c.handlers = append(c.handlers, h)
	}
}

func WithAddr(addr string) Option {
	return func(c *HttpServiceComponent) {
		c.addr = addr
	}
}

func (c *HttpServiceComponent) CommandInitialize(cmd *cobra.Command) {
	for _, h := range c.handlers {
		if ci, ok := h.(svc.CommandInitializer); ok {
			ci.CommandInitialize(cmd)
		}
	}
	cmd.Flags().StringVarP(&c.addr, "http-addr", "", ":8080",
		"Network address and port to list for HTTP requests on.")
}

func (c *HttpServiceComponent) Start(ctx context.Context) error {
	for _, h := range c.handlers {
		if co, ok := h.(svc.Component); ok {
			if err := co.Start(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *HttpServiceComponent) Stop() error {
	for _, h := range c.handlers {
		if co, ok := h.(svc.Component); ok {
			if err := co.Stop(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *HttpServiceComponent) Kill() error {
	for _, h := range c.handlers {
		if co, ok := h.(svc.Component); ok {
			if err := co.Kill(); err != nil {
				return err
			}
		}
	}
	return nil
}
