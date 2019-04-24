package neo

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestAddress2ScriptHash(t *testing.T) {
	hash, err := Address2ScriptHash("AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB")
	assert.NoError(t,err)
	assert.Len(t, hash, 20*2)
	assert.Equal(t, "05a0a304bac8edf51064a9165670cb39fb87439e", string(hash))
}

func TestFormatBigFloat(t *testing.T) {
	v1, err := FormatBigFloat("100000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1", v1)

	v2, err := FormatBigFloat("110000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1.1", v2)

	v3, err := FormatBigFloat("100000000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1000", v3)

	v4, err := FormatBigFloat("10000000000000000000000000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "100000000000000000000", v4)

	v6, err := FormatBigFloat("123456789", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1.23456789", v6)

	v7, err := FormatBigFloat("2386f26fc10000", 16, 8)
	assert.NoError(t, err)
	assert.Equal(t, "100000000", v7)

	v8, err := FormatBigFloat(ReverseBigLitterEndian("0000c16ff28623"), 16, 8)
	assert.NoError(t, err)
	assert.Equal(t, "100000000", v8)
}