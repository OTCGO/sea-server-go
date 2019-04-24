package neo

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestB58encoding(t *testing.T) {
	checkB58Encoding(t, []byte("1"))

	checkB58Encoding(t, []byte("hello world"))
}

func checkB58Encoding(t *testing.T, data []byte) {
	e := b58encode(data)
	d, err := b58decode(e)
	assert.NoError(t, err)
	assert.Equal(t, data, d)
}

func TestBase58CheckEncoding(t *testing.T) {
	checkB58CheckEncoding(t, []byte("1"))

	checkB58CheckEncoding(t, []byte("hello world"))
}

func checkB58CheckEncoding(t *testing.T, data []byte) {
	e := Base58CheckEncode(0, data)
	v, d, err := Base58CheckDecode(e)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint8(0), v)
	assert.Equal(t, data, d)
}
