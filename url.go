package youtube

import (
	"errors"
	"fmt"
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fasthttp"
	"regexp"
)

var RegexStreamID = regexp.MustCompile(`(?i)([a-z0-9_-]{11})`)

type StreamID string

func ExtractStreamID(url string) (StreamID, error) {
	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)

	uri.Parse(nil, bytesutil.Slice(url))

	matches := RegexStreamID.FindSubmatch(uri.RequestURI())
	if matches == nil {
		return "", errors.New("could not find stream id")
	}
	return StreamID(matches[1]), nil
}

func (v StreamID) Valid() error {
	if !RegexStreamID.MatchString(string(v)) {
		return fmt.Errorf("stream id %q is invalid", v)
	}
	return nil
}
