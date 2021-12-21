package measure

import (
	"fmt"
	"strings"
)

var _ Value = (*Metric)(nil)

type Metric struct {
	Name  string
	Value Value
}

func (m *Metric) Encode() string {
	return fmt.Sprintf("%s/%s/%s", m.Value.Type(), m.Name, m.Value.Encode())
}

func (m *Metric) Decode(s string) error {
	chunks := strings.Split(s, "/")
	if len(chunks) != 3 {
		return fmt.Errorf("invalid metric '%s': expected format type/name/value", s)
	}

	if len(chunks[1]) == 0 {
		return fmt.Errorf("invalid metric '%s': name is empty", s)
	}
	m.Name = chunks[1]

	var err error
	if m.Value, err = Type(chunks[0]).New(); err != nil {
		return fmt.Errorf("can not create typed value (%s): %v", chunks[0], err)
	}
	if err = m.Value.Decode(chunks[2]); err != nil {
		return fmt.Errorf("can not read value '%s': %v", chunks[2], err)
	}

	return nil
}

func (*Metric) Type() Type {
	return MetricType
}
