package syncblk

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestSplitHeight(t *testing.T) {
	assert.Len(t, splitHeight(1, 1, 1, 1), 0)

	h1 := splitHeight(1, 2, 1, 1)
	assert.Len(t, h1, 1)
	assert.Equal(t, h1[0].GetInt64("start"), int64(1))
	assert.Equal(t, h1[0].GetInt64("end"), int64(2))

	h2 := splitHeight(1, 12, 2, 5)
	assert.Len(t, h2, 2)
	assert.Equal(t, h2[0].GetInt64("start"), int64(1))
	assert.Equal(t, h2[0].GetInt64("end"), int64(6))
	assert.Equal(t, h2[1].GetInt64("start"), int64(6))
	assert.Equal(t, h2[1].GetInt64("end"), int64(11))

	h3 := splitHeight(1, 12, 4, 5)
	assert.Len(t, h3, 3)
	assert.Equal(t, h3[0].GetInt64("start"), int64(1))
	assert.Equal(t, h3[0].GetInt64("end"), int64(6))
	assert.Equal(t, h3[1].GetInt64("start"), int64(6))
	assert.Equal(t, h3[1].GetInt64("end"), int64(11))
	assert.Equal(t, h3[2].GetInt64("start"), int64(11))
	assert.Equal(t, h3[2].GetInt64("end"), int64(12))
}
