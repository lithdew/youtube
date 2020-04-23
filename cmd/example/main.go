package main

import (
	"fmt"
	"github.com/lithdew/nicehttp"
	"github.com/lithdew/youtube"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Search for the song Animus Vox by The Glitch Mob.

	results, err := youtube.Search("animus vox", 0)
	check(err)

	fmt.Printf("Got %d search result(s).\n\n", results.Hits)

	if len(results.Items) == 0 {
		check(fmt.Errorf("got zero search results"))
	}

	// Get the first search result and print out its details.

	details := results.Items[0]

	fmt.Printf(
		"ID: %q\n\nTitle: %q\nAuthor: %q\nDuration: %q\n\nView Count: %q\nLikes: %d\nDislikes: %d\n\n",
		details.ID,
		details.Title,
		details.Author,
		details.Duration,
		details.Views,
		details.Likes,
		details.Dislikes,
	)

	// Instantiate a player for the first search result.

	player, err := youtube.Load(details.ID)
	check(err)

	// Fetch audio-only direct link.

	stream, ok := player.SourceFormats().AudioOnly().BestAudio()
	if !ok {
		check(fmt.Errorf("no audio-only stream available"))
	}

	audioOnlyFilename := "audio." + stream.FileExtension()

	audioOnlyURL, err := player.ResolveURL(stream)
	check(err)

	fmt.Printf("Audio-only direct link: %q\n", audioOnlyURL)

	// Fetch video-only direct link.

	stream, ok = player.SourceFormats().VideoOnly().BestVideo()
	if !ok {
		check(fmt.Errorf("no video-only stream available"))
	}

	videoOnlyFilename := "video." + stream.FileExtension()

	videoOnlyURL, err := player.ResolveURL(stream)
	check(err)

	fmt.Printf("Video-only direct link: %q\n", videoOnlyURL)

	// Fetch muxed video/audio direct link.

	stream, ok = player.MuxedFormats().BestVideo()
	if !ok {
		check(fmt.Errorf("no muxed stream available"))
	}

	muxedFilename := "muxed." + stream.FileExtension()

	muxedURL, err := player.ResolveURL(stream)
	check(err)

	fmt.Printf("Muxed (video/audio) direct link: %q\n", muxedURL)

	// Download all the links.

	check(nicehttp.DownloadFile(audioOnlyFilename, audioOnlyURL))
	check(nicehttp.DownloadFile(videoOnlyFilename, videoOnlyURL))
	check(nicehttp.DownloadFile(muxedFilename, muxedURL))
}
