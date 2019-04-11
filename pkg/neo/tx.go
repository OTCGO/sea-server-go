package neo

import (
	"github.com/hzxiao/goutil"
	"strings"
)

const (
	MinerTransaction      byte = 0x00
	IssueTransaction      byte = 0x01
	ClaimTransaction      byte = 0x02
	EnrollmentTransaction byte = 0x20
	RegisterTransaction   byte = 0x40
	ContractTransaction   byte = 0x80
	PublishTransaction    byte = 0xd0
	InvocationTransaction byte = 0xd1
)

func IsRegisterNep5AssetTx(tx goutil.Map) bool {
	if tx == nil {
		return false
	}

	if tx.GetString("type") == "InvocationTransaction" &&
		tx.GetInt64("sys_fee") >= 490 &&
		strings.HasSuffix(tx.GetString("script"), "68134e656f2e436f6e74726163742e437265617465") {
			return true
	}
	return false
}

