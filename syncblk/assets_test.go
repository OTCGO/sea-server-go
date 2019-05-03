package syncblk

import (
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestParseNep5Assets(t *testing.T)  {
	var block goutil.Map
	err := neo.Rpc(`https://test4.cityofzion.io:443`, neo.MethodGetBlock, []int{1802320, 1}, &block)
	assert.NoError(t, err)

	asset, err := parseNep5Asset(block.GetMapP("tx/1"))
	assert.NoError(t, err)
	t.Log(goutil.Struct2Json(asset))
}