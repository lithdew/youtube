package youtube

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadPlaylist(t *testing.T) {
	client := NewClient()

	res, err := client.LoadPlaylist("PL25785A39039615CF", 0)
	require.NoError(t, err)

	spew.Dump(res)
}
