package metric

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

type Counter int64

var _ Value = (*Counter)(nil)

func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (c *Counter) LoggerCtx(ctx zerolog.Context) zerolog.Context {
	return logging.UpdateLogCtxWith(ctx.Int64(logging.MetricValueKey, int64(*c)), c.Type())
}

func (c *Counter) Parse(s string) error {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse counter from '%s': %v", s, err)
	}
	*c = Counter(val)
	return nil
}

func (*Counter) Type() Type {
	return CounterType
}
