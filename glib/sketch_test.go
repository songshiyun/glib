package glib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext2Power(t *testing.T) {
	a := assert.New(t)
	a.Equal(int64(16), Next2Power(10))
	a.Equal(int64(1), Next2Power(1))
	a.Equal(int64(2), Next2Power(2))
	a.Equal(int64(4), Next2Power(3))
}

func TestNewCmSketch(t *testing.T) {
	cm := NewCmSketch(128)
	for i := 0 ; i < 1000000; i++ {
		cm.Increment(uint64(i))
	}
	for i := 0 ; i < 100; i++ {
		t.Log(i,cm.Estimate(uint64(i)))
	}
}