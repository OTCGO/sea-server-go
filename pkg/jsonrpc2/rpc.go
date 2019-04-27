package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var jsonNull = json.RawMessage("null")

type JRpcRequest struct {
	ID     int              `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params,omitempty"`
}

func (r JRpcRequest) MarshalJSON() ([]byte, error) {
	r2 := struct {
		Method  string           `json:"method"`
		Params  *json.RawMessage `json:"params,omitempty"`
		ID      int              `json:"id"`
		JSONRPC string           `json:"jsonrpc"`
	}{
		ID:      r.ID,
		Method:  r.Method,
		Params:  r.Params,
		JSONRPC: "2.0",
	}
	return json.Marshal(r2)
}

func (r *JRpcRequest) UnmarshalJSON(data []byte) error {
	var r2 struct {
		Method string           `json:"method"`
		Params *json.RawMessage `json:"params,omitempty"`
		ID     int              `json:"id"`
	}

	// Detect if the "params" field is JSON "null" or just not present
	// by seeing if the field gets overwritten to nil.
	r2.Params = &json.RawMessage{}

	if err := json.Unmarshal(data, &r2); err != nil {
		return err
	}
	r.Method = r2.Method
	if r2.Params == nil {
		r.Params = &jsonNull
	} else if len(*r2.Params) == 0 {
		r.Params = nil
	} else {
		r.Params = r2.Params
	}
	return nil
}

// SetParams sets r.Params to the JSON representation of v. If JSON
// marshaling fails, it returns an error.
func (r *JRpcRequest) SetParams(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.Params = (*json.RawMessage)(&b)
	return nil
}

type JRpcResponse struct {
	ID     int              `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *Error           `json:"error,omitempty"`
}

func (r JRpcResponse) MarshalJSON() ([]byte, error) {
	if (r.Result == nil || len(*r.Result) == 0) && r.Error == nil {
		return nil, errors.New("can't marshal *JRpcResponse (must have result or error)")
	}
	type tmpType JRpcResponse // avoid infinite MarshalJSON recursion
	b, err := json.Marshal(tmpType(r))
	if err != nil {
		return nil, err
	}
	b = append(b[:len(b)-1], []byte(`,"jsonrpc":"2.0"}`)...)
	return b, nil
}

func (r *JRpcResponse) UnmarshalJSON(data []byte) error {
	type tmpType JRpcResponse

	// Detect if the "result" field is JSON "null" or just not present
	// by seeing if the field gets overwritten to nil.
	*r = JRpcResponse{Result: &json.RawMessage{}}

	if err := json.Unmarshal(data, (*tmpType)(r)); err != nil {
		return err
	}
	if r.Result == nil { // JSON "null"
		r.Result = &jsonNull
	} else if len(*r.Result) == 0 {
		r.Result = nil
	}
	return nil
}

func (r *JRpcResponse) UnmarshalResult(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	if r.Result == nil || r.Result == &jsonNull {
		return errors.New("response result is nil")
	}
	
	if err := json.Unmarshal(*(r.Result), &v); err != nil {
		return err
	}
	
	return nil
}

type Error struct {
	Code    int64            `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// Error implements the Go error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc2: code %v message: %s", e.Code, e.Message)
}

var client = &http.Client{
	Timeout: time.Second * 20,
}

func Send(uri string, r *JRpcRequest, result interface{}) error {
	if r == nil {
		return errors.New("doJsonRpc: json-rpc request is nil")
	}

	reqData, err := r.MarshalJSON()
	if err != nil {
		return fmt.Errorf("doJsonRpc: marshal request err(%v)", err)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("doJsonRpc: send http request err(%v)", err)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("doJsonRpc: read http response err(%v)", err)
	}
	defer resp.Body.Close()
	jRpcResp := &JRpcResponse{}
	err = jRpcResp.UnmarshalJSON(respData)
	if err != nil {
		return fmt.Errorf("doJsonRpc: unmarshal json-rpc reponse err(%v)", err)
	}

	return jRpcResp.UnmarshalResult(&result)
}

func SendTimeout(uri string, r *JRpcRequest, timeout time.Duration, result interface{}) error {
	if r == nil {
		return errors.New("doJsonRpc: json-rpc request is nil")
	}

	reqData, err := r.MarshalJSON()
	if err != nil {
		return fmt.Errorf("doJsonRpc: marshal request err(%v)", err)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: timeout,
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("doJsonRpc: send http request err(%v)", err)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("doJsonRpc: read http response err(%v)", err)
	}
	defer resp.Body.Close()
	jRpcResp := &JRpcResponse{}
	err = jRpcResp.UnmarshalJSON(respData)
	if err != nil {
		return fmt.Errorf("doJsonRpc: unmarshal json-rpc reponse err(%v)", err)
	}

	return jRpcResp.UnmarshalResult(&result)
}