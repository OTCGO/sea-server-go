package jsonrpc2

import (
	"encoding/json"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJRpcRequest_MarshalJSON(t *testing.T) {
	r1 := new(JRpcRequest)
	b1, err := r1.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"method":"","id":0,"jsonrpc":"2.0"}`, string(b1))

	r2 := &JRpcRequest{
		ID:     1,
		Method: "run",
	}
	b2, err := r2.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"method":"run","id":1,"jsonrpc":"2.0"}`, string(b2))
}

func TestJRpcRequest_SetParams(t *testing.T) {
	var tables = []struct {
		params interface{}
		result string
	}{
		{params: "string", result: `"string"`},
		{params: 0, result: `0`},
		{params: []int{1, 2, 3}, result: `[1,2,3]`},
		{params: []interface{}{"string", 1}, result: `["string",1]`},
		{
			params: map[string]interface{}{"s": "string", "i": 1},
			result: `{"i":1,"s":"string"}`,
		},
		{params: nil, result: `null`},
	}

	var err error
	for i := range tables {
		r := new(JRpcRequest)
		err = r.SetParams(tables[i].params)
		assert.NoError(t, err)
		assert.Equal(t, tables[i].result, string(*r.Params))
	}
}

func TestJRpcRequest_UnmarshalJSON(t *testing.T) {

}

func TestJRpcResponse_MarshalJSON(t *testing.T) {

}

func TestJRpcResponse_UnmarshalJSON(t *testing.T) {
	var raw = func(msg string) *json.RawMessage {
		var m = json.RawMessage(msg)
		return &m
	}
	var tables = []struct {
		data          string
		resp          JRpcResponse
		result        interface{}
		resultErrFlag bool
	}{
		{
			data:   `{"id":1,"result":1}`,
			resp:   JRpcResponse{ID: 1, Result: raw(`1`)},
			result: 1,
		},
		{
			data:   `{"id":1,"result":"string"}`,
			resp:   JRpcResponse{ID: 1, Result: raw(`"string"`)},
			result: "string",
		},
		{
			data:   `{"id":1,"result":["string",1.1]}`,
			resp:   JRpcResponse{ID: 1, Result: raw(`["string",1.1]`)},
			result: []interface{}{"string", 1.1},
		},
		{
			data:   `{"id":1,"result":{"f":1.1,"s":"string"}}`,
			resp:   JRpcResponse{ID: 1, Result: raw(`{"f":1.1,"s":"string"}`)},
			result: map[string]interface{}{"f": 1.1, "s": "string"},
		},
		{
			data:   `{"id":1,"result":{"f":1.1,"s":"string"}}`,
			resp:   JRpcResponse{ID: 1, Result: raw(`{"f":1.1,"s":"string"}`)},
			result: map[string]interface{}{"f": 1.1, "s": "string"},
		},
		{
			data:          `{"id":1,"result":null}`,
			resp:          JRpcResponse{ID: 1, Result: &jsonNull},
			result:        nil,
			resultErrFlag: true,
		},
		{
			data:          `{"id":1}`,
			resp:          JRpcResponse{ID: 1, Result: nil},
			result:        nil,
			resultErrFlag: true,
		},
		{
			data:          `{"id":1,"error":{"code":1,"message":"string"}}`,
			resp:          JRpcResponse{ID: 1, Result: nil},
			result:        nil,
			resultErrFlag: true,
		},
	}

	var err error
	for i := range tables {
		r := new(JRpcResponse)
		err = r.UnmarshalJSON([]byte(tables[i].data))
		assert.NoError(t, err)
		assert.Equal(t, tables[i].resp.ID, r.ID)
		if tables[i].resp.Result == nil {
			assert.Nil(t, r.Result)
		} else {
			assert.Equal(t, *tables[i].resp.Result, *r.Result)
		}

		var result interface{}
		err = r.UnmarshalResult(&result)
		if tables[i].resultErrFlag {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)

		ty := reflect.TypeOf(result)
		switch ty.Kind() {
		case reflect.Float64:
			assert.Equal(t, goutil.Float64(tables[i].result), goutil.Float64(result))
		default:
			assert.Equal(t, tables[i].result, result)
		}
	}
}

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqData, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		jReq := new(JRpcRequest)
		err = jReq.UnmarshalJSON(reqData)
		assert.NoError(t, err)

		jResp := &JRpcResponse{
			ID:     jReq.ID,
			Result: jReq.Params,
		}
		respData, err := jResp.MarshalJSON()
		assert.NoError(t, err)

		w.Write(respData)
	}))

	req := &JRpcRequest{ID: 1, Method: "test"}
	err := req.SetParams([]string{"s"})
	assert.NoError(t, err)

	var result []string
	err = Send(server.URL, req, &result)
	assert.NoError(t, err)
	assert.Equal(t, []string{"s"}, result)
}
