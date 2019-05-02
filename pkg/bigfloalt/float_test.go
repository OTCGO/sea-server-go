package bigfloalt

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	v1, err := Format("100000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1", v1)

	v2, err := Format("110000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1.1", v2)

	v3, err := Format("100000000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1000", v3)

	v4, err := Format("10000000000000000000000000000", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "100000000000000000000", v4)

	v6, err := Format("123456789", 10, 8)
	assert.NoError(t, err)
	assert.Equal(t, "1.23456789", v6)

	v7, err := Format("2386f26fc10000", 16, 8)
	assert.NoError(t, err)
	assert.Equal(t, "100000000", v7)
}

func TestAdd(t *testing.T) {
	s1, err := Add("0", "0")
	assert.NoError(t, err)
	assert.Equal(t, "0", s1)

	s2, err := Add("1.1", "1.1")
	assert.NoError(t, err)
	assert.Equal(t, "2.2", s2)

	s3, err := Add("100", "200")
	assert.NoError(t, err)
	assert.Equal(t, "300", s3)
}