package metric

import "fmt"

type Value interface {
	fmt.Stringer
	Type() Type
	Parse(string) error
}
