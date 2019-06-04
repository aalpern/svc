package components

import (
	"context"
	"time"

	"github.com/rcrowley/go-metrics"
)

const (
	DefaultMemstatsFrequency = time.Second * 5
	DefaultGCStatsFrequency  = time.Second * 5
)

type RuntimeMetricsComponent struct {
	MemstatsFrequency time.Duration
	GCStatsFrequency  time.Duration
}

func (r *RuntimeMetricsComponent) Start(ctx context.Context) error {
	if r.MemstatsFrequency == 0 {
		r.MemstatsFrequency = DefaultMemstatsFrequency
	}
	if r.GCStatsFrequency == 0 {
		r.GCStatsFrequency = DefaultGCStatsFrequency
	}

	reg := metrics.DefaultRegistry
	metrics.RegisterDebugGCStats(reg)
	go metrics.CaptureDebugGCStats(reg, r.GCStatsFrequency)

	metrics.RegisterRuntimeMemStats(reg)
	go metrics.CaptureRuntimeMemStats(reg, r.MemstatsFrequency)

	return nil
}

func (r *RuntimeMetricsComponent) Stop() error {
	return nil
}

func (r *RuntimeMetricsComponent) Kill() error {
	return nil
}
