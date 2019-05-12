package sync

import (
	"expvar"
	"github.com/OTCGO/sea-server-go/job/node"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func initSuperNode()  {
	superNode = &node.NodeInfo{
		FastestNode:     expvar.NewString("fastestNode"),
		SupportLogNode:  expvar.NewString("supportLogNode"),
	}
	superNode.FastestNode.Set(`https://test3.cityofzion.io:443`)
	superNode.SupportLogNode.Set(`https://test3.cityofzion.io:443`)
}

func TestRpcGetGlobalBalance(t *testing.T) {
	initSuperNode()

	b, err := rpcGetGlobalBalance("e13440dccae716e16fc01adb3c96169d2d08d16581cad0ced0b4e193c472eac1",
		"AGHdThQFJs5kixWuXkgRsbNKz2LrDYDaQB")
	assert.NoError(t, err)
	assert.NotEqual(t, "", b)
}
