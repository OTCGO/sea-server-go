package neo

import (
	"encoding/binary"
	"fmt"
	"github.com/OTCGO/sea-server-go/pkg/neo/opcode"
)

var paramTypes = map[byte]string{
	0x00: "Signature",
	0x01: "Boolean",
	0x02: "Integer",
	0x03: "Hash160",
	0x04: "Hash256",
	0x05: "bytearray",
	0x06: "PublicKey",
	0x07: "String",
	0x10: "Array",
	0xf0: "InteropInterface",
	0xff: "Void",
}

type ContractInfo struct {
	Script       []byte
	Contract     string
	ContractName string
	Version      string
	Parameter    []string
	ReturnType   string
	UseStorage   bool
	DynamicCall  bool
	Author       string
	Email        string
	Description  string
}

func (c *ContractInfo) parse() error {
	return nil
}

func (c *ContractInfo) parseStorageDynamic() error {
	return nil
}

func (c *ContractInfo) parseReturnType() error {
	return nil
}

func (c *ContractInfo) parseParameter() error {
	return nil
}

func ParseNep5Asset(script []byte) (*ContractInfo, error) {
	return nil, nil
}

func parseScriptElem(script []byte) (elem, rest []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
	}()
	var elemLen int
	var start = 2
	mark := script[0]
	switch {
	case mark <= opcode.PUSHBYTES75:
		elemLen = int(mark)
	case mark == opcode.PUSHDATA1:
		elemLen = int(binary.BigEndian.Uint32(script[start:4]))
		start = 4
	case mark == opcode.PUSHDATA2:
		elemLen = int(binary.BigEndian.Uint32(script[start:6]))
		start = 6
	case mark == opcode.PUSHDATA4:
		elemLen = int(binary.BigEndian.Uint32(script[start:10]))
		start = 10
	case mark == opcode.PUSHM1:
		elem = []byte{0xff, 0xff, 0xff, 0xff}
	case mark >= opcode.PUSH1 && mark <= opcode.PUSH16:
		elem = []byte{mark - opcode.PUSH1}
	default:
		return nil, script, nil
	}

	if elemLen == 0 {
		rest = script[start:]
	} else {
		elem = script[start : start+elemLen*2]
		rest = script[start+elemLen*2:]
	}
	return
}
