package neo

import (
	"github.com/OTCGO/sea-server-go/pkg/jsonrpc2"
	"time"
)

const (
	MethodGetBlockCount     = "getblockcount"
	MethodGetBlock          = "getblock"
	MethodInvokeFunction    = "invokefunction"
	MethodGetRawTransaction = "getrawtransaction"
	MethodGetAccountState   = "getaccountstate"
	MethodGetApplicationLog = "getapplicationlog"
)

func Rpc(url, method string, params interface{}, result interface{}) error {
	r := &jsonrpc2.JRpcRequest{
		ID:     1,
		Method: method,
	}
	err := r.SetParams(params)
	if err != nil {
		return err
	}

	return jsonrpc2.Send(url, r, &result)
}

func RpcTimeout(url string, method string, params interface{}, timeout time.Duration, result interface{}) error {
	r := &jsonrpc2.JRpcRequest{
		ID:     1,
		Method: method,
	}
	err := r.SetParams(params)
	if err != nil {
		return err
	}

	return jsonrpc2.SendTimeout(url, r, timeout, &result)
}