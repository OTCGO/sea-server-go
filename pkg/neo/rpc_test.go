package neo

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func init() {
	NeoURI = func() string {
		return `https://test4.cityofzion.io:443`
	}
}

func TestRpc(t *testing.T) {
	var count int
	err := Rpc(MethodGetBlockCount, []interface{}{}, &count)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, count)

	var result goutil.Map
	err = Rpc(MethodGetBlock, []int{count - 1, 1}, &result)
	assert.NoError(t, err)

	t.Log(result)
}
