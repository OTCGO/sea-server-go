package neo

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/pkg/neo/opcode"
	"github.com/hzxiao/goutil"
	"strconv"
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
	Script       string
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

func (c *ContractInfo) parse() (err error) {
	c.Description, err = c.parseString()
	if err != nil {
		return fmt.Errorf("parse description fail: %v", err)
	}
	c.Email, err = c.parseString()
	if err != nil {
		return fmt.Errorf("parse email fail: %v", err)
	}
	c.Author, err = c.parseString()
	if err != nil {
		return fmt.Errorf("parse author fail: %v", err)
	}
	c.Version, err = c.parseString()
	if err != nil {
		return fmt.Errorf("parse version fail: %v", err)
	}
	c.ContractName, err = c.parseString()
	if err != nil {
		return fmt.Errorf("parse name fail: %v", err)
	}
	err = c.parseStorageDynamic()
	if err != nil {
		return err
	}
	err = c.parseReturnType()
	if err != nil {
		return err
	}
	err = c.parseParameter()
	if err != nil {
		return err
	}
	script, err := c.parseString()
	if err != nil {
		return fmt.Errorf("parse contract script fail: %v", err)
	}
	c.Contract = ScriptToHash([]byte(script))

	return nil
}

func (c *ContractInfo) parseStorageDynamic() (err error) {
	var elem interface{}
	elem, c.Script, err = parseScriptElem(c.Script)
	if err != nil {
		return fmt.Errorf("parse storage dynamic err: %v", err)
	}
	if elem == nil {
		return
	}
	v, ok := elem.(int)
	if !ok {
		return
	}

	c.UseStorage = v&0x01 == 0x01
	c.DynamicCall = v&0x02 == 0x02
	return nil
}

func (c *ContractInfo) parseReturnType() (err error) {
	var elem interface{}
	elem, c.Script, err = parseScriptElem(c.Script)
	if err != nil {
		return fmt.Errorf("parse return type err: %v", err)
	}
	if elem == nil {
		return
	}

	str, ok := elem.(string)
	if ok {
		mark, err := strconv.ParseInt(ReverseBigLitterEndian(str), 16, 64)
		if err != nil {
			return fmt.Errorf("parse return type parse str to int err: %v", err)
		}
		c.ReturnType = paramTypes[byte(mark)]
		return err
	}
	iv, ok := elem.(int)
	if ok {
		c.ReturnType = paramTypes[byte(iv)]
	}
	return
}

func (c *ContractInfo) parseParameter() (err error) {
	var elem interface{}
	elem, c.Script, err = parseScriptElem(c.Script)
	if err != nil {
		return fmt.Errorf("parse parameter err: %v", err)
	}
	if elem == nil {
		return
	}

	str, ok := elem.(string)
	if !ok {
		return fmt.Errorf("wrong parameter type: %v", str)
	}

	for i := 0; i < len(str)-1; i = i + 2 {
		mark, err := strconv.ParseInt(str[i:i+2], 16, 64)
		if err != nil {
			return fmt.Errorf("parse parameter str to int: %v", err)
		}
		c.Parameter = append(c.Parameter, paramTypes[byte(mark)])
	}
	return nil
}

func (c *ContractInfo) parseString() (value string, err error) {
	var elem interface{}
	elem, c.Script, err = parseScriptElem(c.Script)
	if err != nil {
		return "", fmt.Errorf("parse string err: %v", err)
	}
	if elem == nil {
		return
	}

	return HexDecode(goutil.String(elem)), nil
}

func parseScriptElem(script string) (elem interface{}, rest string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
	}()
	var elemLen int64
	var start int64 = 2
	mark := []byte(HexDecode(script))[0]
	switch {
	case mark <= opcode.PUSHBYTES75:
		elemLen = int64(mark)
	case mark == opcode.PUSHDATA1:
		elemLen, err = strconv.ParseInt(script[start:4], 16, 16)
		if err != nil {
			return
		}
		start = 4
	case mark == opcode.PUSHDATA2:
		elemLen, err = strconv.ParseInt(ReverseBigLitterEndian(script[start:6]), 16, 32)
		if err != nil {
			return
		}
		start = 6
	case mark == opcode.PUSHDATA4:
		elemLen, err = strconv.ParseInt(ReverseBigLitterEndian(script[start:10]), 16, 64)
		if err != nil {
			return
		}
		start = 10
	case mark == opcode.PUSHM1:
		elem = -1
	case mark >= opcode.PUSH1 && mark <= opcode.PUSH16:
		elem = int(mark - 0x50)
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

func ParseContract(script string) (*ContractInfo, error) {
	contract := &ContractInfo{Script: script}
	err := contract.parse()
	if err != nil {
		return nil, err
	}
	return contract, nil
}
