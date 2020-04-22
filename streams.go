package youtube

import (
	"github.com/valyala/fastjson"
)

type Streams struct {
	id StreamID
	v  *fastjson.Value
}

// ID returns the unique ID pertaining to this stream.
func (s Streams) ID() StreamID {
	return s.id
}

// SourceFormats returns streaming formats that either video-only or audio-only at their highest quality. See the
// documentation for MuxedFormats for lower quality, premuxed streaming formats..
func (s Streams) SourceFormats() Formats {
	fs := s.v.GetArray("streamingData", "adaptiveFormats")
	if len(fs) == 0 {
		return nil
	}

	formats := make(Formats, len(fs))
	for i := range formats {
		formats[i] = ParseFormatJSON(fs[i])
	}
	return formats
}

// MuxedFormats returns premuxed (video/audio-combined) streaming formats. Premuxing comes at an expense of poorer
// video/audio quality.
func (s Streams) MuxedFormats() Formats {
	fs := s.v.GetArray("streamingData", "formats")
	if len(fs) == 0 {
		return nil
	}

	formats := make(Formats, len(fs))
	for i := range formats {
		formats[i] = ParseFormatJSON(fs[i])
	}
	return formats
}

func (s Streams) Title() string {
	return string(s.v.GetStringBytes("videoDetails", "title"))
}

func (s Streams) Author() string {
	return string(s.v.GetStringBytes("videoDetails", "author"))
}

func (s Streams) ChannelID() string {
	return string(s.v.GetStringBytes("videoDetails", "channelId"))
}

func (s Streams) ShortDescription() string {
	return string(s.v.GetStringBytes("videoDetails", "shortDescription"))
}

func (s Streams) Keywords() []string {
	vals := s.v.GetArray("videoDetails", "keywords")

	results := make([]string, 0, len(vals))
	for _, v := range vals {
		results = append(results, string(v.GetStringBytes()))
	}

	return results
}

func (s Streams) AverageRating() float64 {
	return s.v.GetFloat64("videoDetails", "averageRating")
}

func (s Streams) ViewCount() string {
	return string(s.v.GetStringBytes("videoDetails", "viewCount"))
}

func (s Streams) ContextParams() string {
	return string(s.v.GetStringBytes("playabilityStatus", "contextParams"))
}

func (s Streams) Status() string {
	return string(s.v.GetStringBytes("playabilityStatus", "status"))
}

func (s Streams) Reason() string {
	return string(s.v.GetStringBytes("playabilityStatus", "reason"))
}

func (s Streams) PlayableInEmbed() bool {
	return s.v.GetBool("playabilityStatus", "playableInEmbed")
}

func (s Streams) ExpiresInSeconds() string {
	return string(s.v.GetStringBytes("streamingData", "expiresInSeconds"))
}
