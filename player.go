package youtube

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	RegexWatchPlayerConfig = regexp.MustCompile(`ytplayer\.config = ({(?:"\w+":(?:.*?))*});`)
	RegexEmbedPlayerConfig = regexp.MustCompile(`yt\.setConfig\({'PLAYER_CONFIG': (.*?)}\)`)
)

type Player struct {
	Transport
	Assets
	Streams
}

func (p Player) ResolveURL(v Format) (string, error) {
	return p.ResolveURLDeadline(v, zeroTime)
}

func (p Player) ResolveURLTimeout(v Format, timeout time.Duration) (string, error) {
	return p.ResolveURLDeadline(v, time.Now().Add(timeout))
}

func (p Player) ResolveURLDeadline(v Format, deadline time.Time) (string, error) {
	if v.URL != nil {
		return *v.URL, nil
	}

	if v.Cipher == nil {
		return "", errors.New("no url could be found")
	}

	script, err := p.Assets.LoadJSDeadline(p, deadline)
	if err != nil {
		return "", fmt.Errorf("failed to load player script: %w", err)
	}

	return v.Cipher.DecodeURL(script)
}
