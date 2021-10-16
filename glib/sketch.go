package glib

import (
	"math/rand"
	"time"
)

/**
count min sketch ，bloom过滤器的一个变种，
通过4bit(最大次数1111=15)保存一个数的访问频率（tinyLfu认为15次的访问次数是一个较大的访问次数）
同时通过一个刷新机制(当总的访问次数达到一定的次数)来将所有的次数减半来保持相对活性，
这样历史(访问量较大的)数据能够逐步被淘汰
*/
const (
	cmDepth = 4
)

type CountMinSketch struct {
	rows [cmDepth]cmRow
	seed [cmDepth]uint64
	mask uint64
}

func NewCmSketch(counter int64) *CountMinSketch {
	if counter == 0 {
		panic("bad counter")
	}
	counter = Next2Power(counter)
	sketch := &CountMinSketch{mask: uint64(counter - 1)}
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < cmDepth; i++ {
		sketch.seed[i] = source.Uint64()
		sketch.rows[i] = newRow(counter)
	}
	return sketch
}

func (cm *CountMinSketch) Increment(data uint64) {
	for i := range cm.rows {
		d := (data ^ cm.seed[i]) & cm.mask
		cm.rows[i].incr(d)
	}
}

func (cm *CountMinSketch) Estimate(data uint64) int64 {
	min := byte(255)
	for i := range cm.rows {
		val := cm.rows[i].get(data)
		if val < min {
			min = val
		}
	}
	return int64(min)
}

func (cm *CountMinSketch) Reset() {
	for _, item := range cm.rows {
		item.reset()
	}
}

func (cm *CountMinSketch) Clear() {
	for _, item := range cm.rows {
		item.clear()
	}
}

type cmRow []byte

func newRow(counter int64) cmRow {
	return make(cmRow, counter/2)
}

func (c cmRow) get(d uint64) byte {
	return (c[d/2] >> ((d & 1) * 4)) & 0x0f
}

func (c cmRow) incr(d uint64) {
	i := d / 2
	s := (d & 1) * 4 // ==0 || 4
	v := (c[i] >> s) & 0x0f
	if v < 15 {
		c[i] += 1 << s
	}
}

func (c cmRow) reset() {
	for i := range c {
		c[i] = (c[i] >> 1) & 0x0f // 减半
	}
}

func (c cmRow) clear() {
	for i := range c {
		c[i] = 0
	}
}

func Next2Power(x int64) int64 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}
