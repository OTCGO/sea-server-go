package httptest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hzxiao/goutil"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var GinEngine *gin.Engine

func Get(url string) (goutil.Map, error) {
	return fetch("GET", url, "", nil)
}

func Post(url, cType string, reader io.Reader) (goutil.Map, error) {
	return fetch("POST", url, cType, reader)
}

func fetch(method, url, cType string, reader io.Reader) (goutil.Map, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if cType != "" {
		req.Header.Set("Content-Type", cType)
	}

	resp, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	return fetchJson(resp)
}

func doRequest(req *http.Request) (http.ResponseWriter, error) {
	if GinEngine == nil {
		return nil, errors.New("gin.engine should be reinitialized")
	}

	w := httptest.NewRecorder()
	GinEngine.ServeHTTP(w, req)
	return w, nil
}

func fetchJson(w http.ResponseWriter) (goutil.Map, error) {
	resp, ok := w.(*httptest.ResponseRecorder)
	if !ok {
		return nil, errors.New("unknown response")
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data goutil.Map
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func PostJSON(url string, body interface{}) (goutil.Map, error) {
	return Post(url, "application/json", bytes.NewBufferString(goutil.Struct2Json(body)))
}
