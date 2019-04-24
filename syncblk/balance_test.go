package syncblk

import (
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func init() {
	neo.NeoURI = func() string {
		return `https://test3.cityofzion.io:443`
	}
}

func TestRpcGetGlobalBalance(t *testing.T) {
	b, err := rpcGetGlobalBalance("e13440dccae716e16fc01adb3c96169d2d08d16581cad0ced0b4e193c472eac1",
		"AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB")
	assert.NoError(t, err)
	assert.NotEqual(t, "", b)
}
