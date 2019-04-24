package neo

import (
	"github.com/OTCGO/sea-server-go/pkg/jsonrpc2"
)

const (
	MethodGetBlockCount     = "getblockcount"
	MethodGetBlock          = "getblock"
	MethodInvokeFunction    = "invokefunction"
	MethodGetRawTransaction = "getrawtransaction"
	MethodGetAccountState   = "getaccountstate"
)

var NeoURI = func() string {
	panic("neo uri isn't be initialized")
}

func Rpc(method string, params interface{}, result interface{}) error {
	r := &jsonrpc2.JRpcRequest{
		ID:     1,
		Method: method,
	}
	err := r.SetParams(params)
	if err != nil {
		return err
	}

	return jsonrpc2.Send(NeoURI(), r, &result)
}
