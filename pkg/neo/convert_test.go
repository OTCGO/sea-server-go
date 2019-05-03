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

func TestScriptHash2Address(t *testing.T) {
	address, err := ScriptHash2Address(HexDecodeString("05a0a304bac8edf51064a9165670cb39fb87439e"))
	assert.NoError(t, err)
	assert.Equal(t, "AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB", address)
}