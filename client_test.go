package youtube

import (
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
	"testing"
)

func TestLoadEmbedPlayer(t *testing.T) {
	client := NewClient()

	p, err := client.LoadEmbedPlayer("pAsDzfbLM8Y")
	require.NoError(t, err)
	require.NotZero(t, p)
}

func BenchmarkURLAppend(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	id := "hello_world"
	offset := uint(12)

	for i := 0; i < b.N; i++ {
		uri := []byte("https://www.youtube.com/list_ajax?style=json&action_get_list=1")

		uri = append(uri, "&list="...)
		uri = append(uri, id...)

		uri = append(uri, "&index="...)
		uri = fasthttp.AppendUint(uri, int(offset))

		uri = append(uri, "&hl="...)
		uri = append(uri, "en"...)
	}
}

func BenchmarkURLValues(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	id := "hello_world"
	offset := uint(12)

	for i := 0; i < b.N; i++ {
		uri := []byte("https://www.youtube.com/list_ajax?")

		uri = append(uri,
			url.Values{
				"style":           {"json"},
				"action_get_list": {strconv.FormatUint(1, 10)},
				"list":            {id},
				"index":           {strconv.FormatUint(uint64(offset), 10)},
				"hl":              {"en"},
			}.Encode()...,
		)
	}
}
