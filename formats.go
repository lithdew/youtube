package youtube

import (
	"sort"
)

var AudioQuality = map[string]int{
	"AUDIO_QUALITY_LOW":    0,
	"AUDIO_QUALITY_MEDIUM": 1,
	"AUDIO_QUALITY_HIGH":   2,
}

var VideoQuality = map[string]int{
	"tiny":   0,
	"low":    1,
	"medium": 2,
	"large":  3,
	"hd1440": 4,
	"hd2160": 5,
}

type Formats []Format

func (f Formats) VideoOnly() Formats {
	return FilterVideoStreams(f)
}

func (f Formats) AudioOnly() Formats {
	return FilterAudioStreams(f)
}

func (f Formats) SortByVideoQuality() Formats {
	return SortByVideoQuality(f)
}

func (f Formats) SortByAudioQuality() Formats {
	return SortByAudioQuality(f)
}

func (f Formats) BestVideo() (Format, bool) {
	return SearchForBestVideoQuality(f)
}

func (f Formats) BestAudio() (Format, bool) {
	return SearchForBestAudioQuality(f)
}

func SearchForBestVideoQuality(formats Formats) (Format, bool) {
	formats = SortByVideoQuality(formats)
	if len(formats) == 0 {
		return Format{}, false
	}

	return formats[0], true
}

func SearchForBestAudioQuality(formats Formats) (Format, bool) {
	formats = SortByAudioQuality(formats)
	if len(formats) == 0 {
		return Format{}, false
	}

	return formats[0], true
}

func FilterAudioStreams(formats Formats) Formats {
	filtered := formats[:0]
	for _, format := range formats {
		if format.AudioQuality == nil {
			continue
		}
		filtered = append(filtered, format)
	}
	return filtered
}

func FilterVideoStreams(formats Formats) Formats {
	filtered := formats[:0]
	for _, format := range formats {
		if format.FPS == nil {
			continue
		}
		filtered = append(filtered, format)
	}
	return filtered
}

func SortByVideoQuality(formats Formats) Formats {
	sort.Slice(formats, func(i, j int) bool {
		if formats[i].Bitrate > formats[j].Bitrate {
			return true
		}
		return VideoQuality[formats[i].Quality] > VideoQuality[formats[j].Quality]
	})

	return formats
}

func SortByAudioQuality(formats Formats) Formats {
	sort.Slice(formats, func(i, j int) bool {
		if formats[i].Bitrate > formats[j].Bitrate {
			return true
		}

		a, b := 0, 0

		if ap := formats[i].AudioQuality; ap != nil {
			a = AudioQuality[*ap]
		}

		if bp := formats[j].AudioQuality; bp != nil {
			b = AudioQuality[*bp]
		}

		return a > b
	})

	return formats
}
