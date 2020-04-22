package youtube

import (
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fastjson"
)

type PlaylistResult struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Views       uint   `json:"views"`

	Items []ListItem `json:"video"`
}

func ParsePlaylistResultJSON(v *fastjson.Value) PlaylistResult {
	vals := v.GetArray("video")

	r := PlaylistResult{
		Title:       bytesutil.String(v.GetStringBytes("title")),
		Author:      bytesutil.String(v.GetStringBytes("author")),
		Description: bytesutil.String(v.GetStringBytes("description")),
		Views:       v.GetUint("views"),
		Items:       make([]ListItem, 0, len(vals)),
	}

	for _, val := range vals {
		r.Items = append(r.Items, ParseListItem(val))
	}

	return r
}
