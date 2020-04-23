# youtube

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lithdew/youtube)
[![Discord Chat](https://img.shields.io/discord/697002823123992617)](https://discord.gg/HZEbkeQ)

A library for retrieving metadata and obtaining direct links to video-only/audio-only/muxed streams of videos on YouTube in Go.

## Inspiration

I was always curious on how YouTube downloader sites and software worked. So, I took the time to trace and reverse-engineer how HTTP requests are made to YouTube to then allow for video/audio from YouTube to be played on mobile/in the browser.

As a result, this library in Go allows for you to retrieve direct links to video-only/audio-only/muxed streams to YouTube videos, retrieve metadata of YouTube videos, and make searches for videos on YouTube without any need for an API key and without any usage quota restrictions.

This library also exposes all of its methods for decoding encrypted direct links to media, and decoding/parsing both JSON and URL-encoded documents returned by YouTube.

Many thanks to the library and blog posts from [Tyrrrz/YoutubeExplode](https://github.com/Tyrrrz/YoutubeExplode) for giving me a good direction on how to get down-and-dirty tracing and decoding HTTP requests on YouTube.

## Features

- Does not use require an API key or have any usage quotas.
- Retrieve direct links to video-only/audio-only/muxed streams to YouTube videos.
- Retrieve metadata of videos or playlists on YouTube.
- Search for videos/audio on YouTube.
- Set timeouts/deadlines for all methods.
- Minimal dependencies.
- Concurrency-safe.

## Setup

```
$ go get github.com/lithdew/youtube
```

## Example

This example uses my library [lithdew/nicehttp](https://github.com/lithdew/nicehttp) for downloading video/audio from YouTube as fast as possible in fixed-sized chunks with multiple workers working in parallel.

It searches for the song `The Glitch Mob - Animus Vox` on YouTube and downloads its audio, video, and muxed versions to disk.

It additionally prints all metadata pertaining to the first video it finds.

```go
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
```

You can run this example by running:

```shell
$ go run github.com/lithdew/youtube/cmd/example
```

An extended example is also provided for downloading to disk the highest quality audio-only stream of any YouTube video. It may be run by running:

```shell
$ go run github.com/lithdew/youtube/cmd/music https://www.youtube.com/watch?v=jPan651rVMs
```

## What's missing?

Although this library is feature-complete for several use cases, there are a few things intentionally missing as I started this library as just a side project. The features missing are:

1. Fetching video captions.
2. Fetching video comments.
3. DASH manifest support.
4. Livestream support.

Should there be enough demand however, I'll look into taking some time to incorporate these features into this library :).

Open a Github issue to express your interest!