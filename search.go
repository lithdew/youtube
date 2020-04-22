package youtube

import (
	"github.com/valyala/fastjson"
)

type SearchResult struct {
	Hits  uint       `json:"hits"`
	Items []ListItem `json:"video"`
}

func ParseSearchResultJSON(v *fastjson.Value) SearchResult {
	vals := v.GetArray("video")

	r := SearchResult{
		Hits:  v.GetUint("hits"),
		Items: make([]ListItem, 0, len(vals)),
	}

	for _, val := range vals {
		r.Items = append(r.Items, ParseListItem(val))
	}

	return r
}
