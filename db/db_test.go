package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func init()  {

}

func TestInitDB(t *testing.T) {
	err := InitDB("root:123456@/mysql?charset=utf8&loc=Local&parseTime=true")
	assert.NoError(t, err)
}