package syncblk

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestFormatAppLog(t *testing.T) {
	log1 := goutil.Map{
		"vmstate": "HALT, BREAK",
		"key":     "1",
	}
	appLog1, ok := formatAppLog(log1)
	assert.True(t, ok)
	assert.Equal(t, "1", appLog1["key"])

	log2 := goutil.Map{
		"vmstate": "FAULT, BREAK",
		"key":     "1",
	}
	appLog2, ok := formatAppLog(log2)
	assert.False(t, ok)
	assert.Nil(t, appLog2)

	log3 := goutil.Map{
		"executions": []goutil.Map{
			{
				"vmstate": "HALT, BREAK",
				"key":     "2",
			},
		},
	}
	appLog3, ok := formatAppLog(log3)
	assert.True(t, ok)
	assert.Equal(t, "2", appLog3["key"])
}

func TestSyncHistory_Sync(t *testing.T) {
	initSuperNode()

	block, err := new(SyncBlock).Block(75)
	assert.NoError(t, err)

	err = new(SyncHistory).Sync(block)
	assert.NoError(t, err)
}

func TestSyncUtxo_Sync(t *testing.T) {
	initSuperNode()

	block, err := new(SyncBlock).Block(82)
	assert.NoError(t, err)

	err = new(SyncUtxo).Sync(block)
	assert.NoError(t, err)
}