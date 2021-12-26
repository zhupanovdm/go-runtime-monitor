package metric

import (
	"fmt"
	"strconv"
)

type Counter int64

var _ Value = (*Counter)(nil)

func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (c *Counter) Parse(s string) error {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse counter from '%s': %v", s, err)
	}
	*c = Counter(val)
	return nil
}

func (Counter) Type() Type {
	return CounterType
}
