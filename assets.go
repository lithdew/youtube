package youtube

import (
	"errors"
	"fmt"
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fastjson"
	"time"
)

type Assets struct {
	CSS string `json:"css"`
	JS  string `json:"js"`
}

func ParseAssetsJSON(v *fastjson.Value) Assets {
	return Assets{
		CSS: string(v.GetStringBytes("css")),
		JS:  string(v.GetStringBytes("js")),
	}
}

func (p Assets) LoadCSS(t Transport) (string, error) {
	return p.LoadCSSDeadline(t, zeroTime)
}

func (p Assets) LoadCSSTimeout(t Transport, timeout time.Duration) (string, error) {
	return p.LoadCSSDeadline(t, time.Now().Add(timeout))
}

func (p Assets) LoadCSSDeadline(t Transport, deadline time.Time) (string, error) {
	if p.CSS == "" {
		return "", errors.New("could not find url to player css")
	}

	buf, err := t.DownloadBytesDeadline(nil, "https://www.youtube.com"+p.CSS, deadline)
	if err != nil {
		return "", fmt.Errorf("failed to download player css: %w", err)
	}

	return bytesutil.String(buf), nil
}

func (p Assets) LoadJS(t Transport) (string, error) {
	return p.LoadJSDeadline(t, zeroTime)
}

func (p Assets) LoadJSTimeout(t Transport, timeout time.Duration) (string, error) {
	return p.LoadJSDeadline(t, time.Now().Add(timeout))
}

func (p Assets) LoadJSDeadline(t Transport, deadline time.Time) (string, error) {
	if p.JS == "" {
		return "", errors.New("could not find url to player script")
	}

	buf, err := t.DownloadBytesDeadline(nil, "https://www.youtube.com"+p.JS, deadline)
	if err != nil {
		return "", fmt.Errorf("failed to download player script: %w", err)
	}

	return bytesutil.String(buf), nil
}
