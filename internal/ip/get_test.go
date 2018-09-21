package ip_test

import (
	"testing"

	"github.com/adamwg/golab-examples/internal/ip"
	"github.com/stretchr/testify/assert"
)

func TestGetSuccess(t *testing.T) {
	ip, err := ip.GetIP()
	assert.NoError(t, err)
	assert.NotNil(t, ip)
}
