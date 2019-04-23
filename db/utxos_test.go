package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestInsertOrUpdateUpt(t *testing.T) {
	deleteAll()

	upt := &Upt{Address: "1", Asset: "2", UpdateHeight: 1}
	ok, err := InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.False(t, ok)

	upt.UpdateHeight = 2
	ok, err = InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestInsertOrUpdateBalance(t *testing.T) {
	deleteAll()

	b := &Balance{
		Address:           "1",
		Asset:             "2",
		Value:             "3",
		LastUpdatedHeight: 2,
	}
	ok, err := InsertOrUpdateBalance(b)
	assert.NoError(t, err)
	assert.True(t, ok)
}
