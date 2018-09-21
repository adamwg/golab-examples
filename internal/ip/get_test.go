package ip_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/adamwg/golab-examples/internal/ip"
	"github.com/stretchr/testify/assert"
)

type failReader struct{}

func (fr *failReader) Read(p []byte) (int, error) {
	return 0, errors.New("failed to read")
}

func TestGetSuccess(t *testing.T) {
	respBody := ioutil.NopCloser(
		bytes.NewBufferString("127.0.0.1\n"),
	)

	resp := &http.Response{
		Body: respBody,
	}

	getter := &ip.MockHTTPGetter{}
	getter.On("Get", "https://icanhazip.com").
		Return(resp, nil).
		Once()

	ip, err := ip.GetIP(getter)
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, net.IPv4(127, 0, 0, 1), ip)
}

func TestGetGarbageResponse(t *testing.T) {
	respBody := ioutil.NopCloser(
		bytes.NewBufferString("this is not an ip address!\n"),
	)

	resp := &http.Response{
		Body: respBody,
	}

	getter := &ip.MockHTTPGetter{}
	getter.On("Get", "https://icanhazip.com").
		Return(resp, nil).
		Once()

	_, err := ip.GetIP(getter)
	assert.Error(t, err)
}

func TestGetReadFailure(t *testing.T) {
	respBody := ioutil.NopCloser(
		&failReader{},
	)

	resp := &http.Response{
		Body: respBody,
	}

	getter := &ip.MockHTTPGetter{}
	getter.On("Get", "https://icanhazip.com").
		Return(resp, nil).
		Once()

	_, err := ip.GetIP(getter)
	assert.Error(t, err)
}

func TestGetRequestFailure(t *testing.T) {
	getter := &ip.MockHTTPGetter{}
	getter.On("Get", "https://icanhazip.com").
		Return(nil, errors.New("request failed")).
		Once()

	_, err := ip.GetIP(getter)
	assert.Error(t, err)
}
