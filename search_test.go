package youtube

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearch(t *testing.T) {
	client := NewClient()

	results, err := client.Search("animus vox", 0)
	require.NoError(t, err)

	require.NotZero(t, results.Hits)
	require.NotEmpty(t, results.Items)
}
