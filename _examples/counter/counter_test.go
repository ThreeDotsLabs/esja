package counter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	counter "postcard"
)

func TestCounter(t *testing.T) {
	c, err := counter.NewCounter("ID")
	require.NoError(t, err)

	require.Equal(t, 0, c.CurrentValue())

	err = c.IncrementBy(10)
	require.NoError(t, err)

	require.Equal(t, 10, c.CurrentValue())

	err = c.IncrementBy(20)
	require.NoError(t, err)
	err = c.IncrementBy(10)
	require.NoError(t, err)

	require.Equal(t, 40, c.CurrentValue())
}
