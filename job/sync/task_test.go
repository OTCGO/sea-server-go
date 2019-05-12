package sync

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestSplitHeight(t *testing.T) {
	h, c := splitHeight(1, 1, 1, 1)
	assert.Len(t, h, 0)
	assert.Equal(t, 0, c)

	h1, c1 := splitHeight(1, 2, 1, 1)
	assert.Len(t, h1, 1)
	assert.Equal(t, 1, c1)
	assert.Equal(t, h1[0]["start"], 1)
	assert.Equal(t, h1[0]["end"], 2)

	h2, c2 := splitHeight(1, 12, 2, 5)
	assert.Equal(t, 10, c2)
	assert.Len(t, h2, 2)
	assert.Equal(t, h2[0]["start"], 1)
	assert.Equal(t, h2[0]["end"], 6)
	assert.Equal(t, h2[1]["start"], 6)
	assert.Equal(t, h2[1]["end"], 11)

	h3, c3 := splitHeight(1, 12, 4, 5)
	assert.Len(t, h3, 3)
	assert.Equal(t, 11, c3)
	assert.Equal(t, h3[0]["start"], 1)
	assert.Equal(t, h3[0]["end"], 6)
	assert.Equal(t, h3[1]["start"], 6)
	assert.Equal(t, h3[1]["end"], 11)
	assert.Equal(t, h3[2]["start"], 11)
	assert.Equal(t, h3[2]["end"], 12)
}
