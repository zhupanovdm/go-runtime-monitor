package metric

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

type Gauge float64

var _ Value = (*Gauge)(nil)

func (g Gauge) String() string {
	return fmt.Sprintf("%.3f", g)
}

func (g *Gauge) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	if g == nil {
		return ctx
	}
	return logging.LogCtxUpdateWith(ctx.Float64(logging.MetricValueKey, float64(*g)), g.Type())
}

func (g *Gauge) Parse(s string) error {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("can't parse gauge from '%s': %v", s, err)
	}
	*g = Gauge(val)
	return nil
}

func (Gauge) Type() Type {
	return GaugeType
}
