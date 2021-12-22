package measure

import (
	"fmt"
	"strconv"
)

type Counter int64

var _ Value = (*Counter)(nil)

func (c Counter) Encode() string {
	return fmt.Sprintf("%d", c)
}

func (c *Counter) Decode(s string) error {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*c = Counter(val)
	return nil
}

func (Counter) Type() Type {
	return CounterType
}
