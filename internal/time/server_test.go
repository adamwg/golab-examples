package time

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adamwg/golab-examples/internal/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	// fakeNow is the reference time used in the time package for format
	// strings.
	fakeNow           = time.Unix(1136239445, 0)
	fakeNowString     = fakeNow.Format(time.RFC3339) + "\n"
	arbitraryTime     = time.Unix(1234567890, 0)
	arbitrarySince    = arbitraryTime.Sub(fakeNow)
	arbitrarySinceStr = arbitrarySince.String() + "\n"

	mockTime = &mocks.Time{}
)

func init() {
	now = mockTime.Now
	since = mockTime.Since
	parseTime = mockTime.Parse
}

func TestServeNowSuccess(t *testing.T) {
	mockTime.On("Now").Return(fakeNow).Once()

	w := httptest.NewRecorder()
	serveNow(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fakeNowString, w.Body.String())

	mockTime.AssertExpectations(t)
}

func TestServeSinceSuccess(t *testing.T) {
	mockTime.On("Parse", time.RFC3339, fakeNowString).
		Return(fakeNow, nil).
		Once()
	mockTime.On("Since", fakeNow).
		Return(arbitrarySince).
		Once()

	w := httptest.NewRecorder()
	serveSince(w, fakeNowString)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, arbitrarySinceStr, w.Body.String())

	mockTime.AssertExpectations(t)
}

func TestServeSinceParseError(t *testing.T) {
	mockTime.On("Parse", time.RFC3339, fakeNowString).
		Return(time.Time{}, errors.New("parse error")).
		Once()

	w := httptest.NewRecorder()
	serveSince(w, fakeNowString)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockTime.AssertExpectations(t)
}
