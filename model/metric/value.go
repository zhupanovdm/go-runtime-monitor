package metric

import (
	"fmt"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/logging"
)

type Value interface {
	fmt.Stringer
	logging.LogCtxProvider
	Type() Type
	Parse(string) error
}
