package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClock_Now(t *testing.T) {
	c := NewClock()
	assert.True(t, time.Now().Before(c.Now()))
}

func TestClock_UTC(t *testing.T) {
	c := NewClock()
	utc, err := time.LoadLocation("UTC")
	require.NoError(t, err)
	assert.True(t, time.Now().Before(c.UTC()))
	assert.Equal(t, utc, c.UTC().Location())
}
