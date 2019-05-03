package neo

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestRpc(t *testing.T) {
	url := `https://test3.cityofzion.io:443`
	var count int
	err := Rpc(url, MethodGetBlockCount, []interface{}{}, &count)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, count)

	var result goutil.Map
	err = Rpc(url, MethodGetBlock, []int{count - 1, 1}, &result)
	assert.NoError(t, err)

	t.Log(result)
}
