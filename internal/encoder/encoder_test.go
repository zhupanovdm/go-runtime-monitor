package encoder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCounter(t *testing.T) {
	c := counter(10)
	cnt := NewCounter("foo", 10)
	if assert.IsType(t, &metric{}, cnt) && assert.Equal(t, "foo", cnt.Type()) {
		m := cnt.(*metric)
		assert.Equal(t, c.Encode(), m.value.Encode())
	}
}

func TestNewGaugeF(t *testing.T) {
	g := gauge(0.5)
	gau := NewGaugeF("foo", 0.5)
	if assert.IsType(t, &metric{}, gau) && assert.Equal(t, "foo", gau.Type()) {
		m := gau.(*metric)
		assert.Equal(t, g.Encode(), m.value.Encode())
	}
}

func TestNewGaugeI(t *testing.T) {
	g := gauge(5)
	gau := NewGaugeI("foo", 5)
	if assert.IsType(t, &metric{}, gau) && assert.Equal(t, "foo", gau.Type()) {
		m := gau.(*metric)
		assert.Equal(t, g.Encode(), m.value.Encode())
	}
}
