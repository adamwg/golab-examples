package time

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServeNowE2E(t *testing.T) {
	mockTime.On("Now").Return(fakeNow).Once()

	srv := httptest.NewServer(NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, fakeNowString, string(body))
}

func TestServeSinceE2E(t *testing.T) {
	mockTime.On("Parse", time.RFC3339, fakeNowString).
		Return(fakeNow, nil).
		Once()
	mockTime.On("Since", fakeNow).
		Return(arbitrarySince).
		Once()

	srv := httptest.NewServer(NewHandler())
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, fakeNowString))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, arbitrarySinceStr, string(body))
}
