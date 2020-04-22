package youtube

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadEmbedPlayer(t *testing.T) {
	client := NewClient()

	p, err := client.LoadEmbedPlayer("pAsDzfbLM8Y")
	require.NoError(t, err)
	require.NotZero(t, p)
}
