package netutil_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"../netutil"
)

func TestWriteResponse(t *testing.T) {
	t.Run("good json body and code", func(t *testing.T) {
		w := httptest.NewRecorder()
		jsonBody := []byte("{\"msg\": \"success\"}")
		netutil.WriteResponse(jsonBody, http.StatusOK, w)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		marshalledJsonBody, _ := json.Marshal(jsonBody)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, body, marshalledJsonBody)
	})
	t.Run("bad json body", func(t *testing.T) {
		w := httptest.NewRecorder()
		badJson := make(chan int)
		netutil.WriteResponse(badJson, http.StatusOK, w)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, body, []byte{})
	})
	t.Run("nil body", func(t *testing.T) {
		w := httptest.NewRecorder()
		netutil.WriteResponse(nil, http.StatusOK, w)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, body, []byte{})
	})
	t.Run("body not allowed for status", func(t *testing.T) {
		w := httptest.NewRecorder()
		netutil.WriteResponse(nil, http.StatusNoContent, w)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, resp.StatusCode, http.StatusNoContent)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, body, []byte{})
	})
}
