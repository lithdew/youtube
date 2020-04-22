package youtube

import (
	"github.com/lithdew/bytesutil"
	"github.com/valyala/fastjson"
)

type Format struct {
	AverageBitrate   uint   `json:"averageBitrate"`
	ApproxDurationMs string `json:"approxDurationMs"`
	ContentLength    string `json:"contentLength"`
	Bitrate          uint   `json:"bitrate"`

	URL    *string `json:"url,omitempty"`
	Cipher *Cipher `json:"cipher,omitempty"`

	Quality      string `json:"quality"`
	QualityLabel string `json:"qualityLabel"`

	ITag     uint   `json:"itag"`
	MIMEType string `json:"mimeType"`

	Width  uint `json:"width"`
	Height uint `json:"height"`

	FPS       *uint      `json:"fps,omitempty"`
	ColorInfo *ColorInfo `json:"colorInfo,omitempty"`

	AudioQuality    *string `json:"audioQuality,omitempty"`
	AudioChannels   *uint   `json:"audioChannels,omitempty"`
	AudioSampleRate *string `json:"audioSampleRate,omitempty"`

	InitRange  *TimeRange `json:"initRange,omitempty"`
	IndexRange *TimeRange `json:"indexRange,omitempty"`

	LastModified    string `json:"lastModified"`
	HighReplication bool   `json:"highReplication,omitempty"`

	ProjectionType string `json:"projectionType"`
}

func (f Format) FileExtension() string {
	return ITags[f.ITag].Extension
}

type ColorInfo struct {
	Primaries               string `json:"primaries"`
	TransferCharacteristics string `json:"transferCharacteristics"`
	MatrixCoefficients      string `json:"matrixCoefficients"`
}

func ParseColorInfoJSON(v *fastjson.Value) ColorInfo {
	return ColorInfo{
		Primaries:               string(v.GetStringBytes("primaries")),
		TransferCharacteristics: string(v.GetStringBytes("transferCharacteristics")),
		MatrixCoefficients:      string(v.GetStringBytes("matrixCoefficients")),
	}
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func ParseTimeRangeJSON(v *fastjson.Value) TimeRange {
	return TimeRange{
		Start: string(v.GetStringBytes("start")),
		End:   string(v.GetStringBytes("end")),
	}
}

func ParseFormatJSON(v *fastjson.Value) Format {
	var format Format

	format.AverageBitrate = v.GetUint("averageBitrate")
	format.ApproxDurationMs = bytesutil.String(v.GetStringBytes("approxDurationMs"))
	format.ContentLength = bytesutil.String(v.GetStringBytes("contentLength"))
	format.Bitrate = v.GetUint("bitrate")

	if u := v.GetStringBytes("url"); len(u) > 0 {
		format.URL = func(s string) *string { return &s }(bytesutil.String(u))
	}

	if u := v.Get("cipher"); u != nil {
		format.Cipher = func(c Cipher) *Cipher { return &c }(ParseCipherJSON(u))
	}

	format.Quality = bytesutil.String(v.GetStringBytes("quality"))
	format.QualityLabel = bytesutil.String(v.GetStringBytes("qualityLabel"))

	format.ITag = v.GetUint("itag")
	format.MIMEType = bytesutil.String(v.GetStringBytes("mimeType"))

	format.Width = v.GetUint("width")
	format.Height = v.GetUint("height")

	format.FPS = func(u uint) *uint { return &u }(v.GetUint("fps"))

	if colorInfo := v.Get("colorInfo"); colorInfo != nil {
		format.ColorInfo = func(i ColorInfo) *ColorInfo { return &i }(ParseColorInfoJSON(colorInfo))
	}

	if audioQuality := v.GetStringBytes("audioQuality"); len(audioQuality) > 0 {
		format.AudioQuality = func(s string) *string { return &s }(string(audioQuality))
	}

	format.AudioChannels = func(u uint) *uint { return &u }(v.GetUint("audioChannels"))

	if audioSampleRate := v.GetStringBytes("audioSampleRate"); len(audioSampleRate) > 0 {
		format.AudioSampleRate = func(s string) *string { return &s }(string(audioSampleRate))
	}

	if initRange := v.Get("initRange"); initRange != nil {
		format.InitRange = func(t TimeRange) *TimeRange { return &t }(ParseTimeRangeJSON(initRange))
	}

	if indexRange := v.Get("indexRange"); indexRange != nil {
		format.IndexRange = func(t TimeRange) *TimeRange { return &t }(ParseTimeRangeJSON(indexRange))
	}

	format.HighReplication = v.GetBool("highReplication")
	format.LastModified = bytesutil.String(v.GetStringBytes("lastModified"))

	format.ProjectionType = bytesutil.String(v.GetStringBytes("projectionType"))

	return format
}
