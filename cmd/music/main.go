package main

import (
	"flag"
	"fmt"
	"github.com/lithdew/nicehttp"
	"github.com/lithdew/youtube"
	"log"
	"path"
	"regexp"
	"strings"
)

var (
	regexSeparators      = regexp.MustCompile(`[ &_=+:]`)
	regexLegalCharacters = regexp.MustCompile(`[^[:alnum:]-.]`)
)

func normalizeFileName(str string) string {
	name := strings.Trim(path.Clean(path.Base(strings.ToLower(str))), " ")

	name = regexSeparators.ReplaceAllString(name, "-")
	name = regexLegalCharacters.ReplaceAllString(name, "")

	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	return name
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	client := youtube.NewClient()

	for _, src := range flag.Args() {
		id, err := youtube.ExtractStreamID(src)
		check(err)

		player, err := client.Load(id)
		check(err)

		fmt.Printf(
			"Title: %q\nAuthor: %q\nView Count: %s\n\n",
			player.Title(),
			player.Author(),
			player.ViewCount(),
		)

		stream, ok := player.SourceFormats().AudioOnly().BestAudio()
		if !ok {
			check(fmt.Errorf("no audio available for video id %q", id))
		}

		url, err := player.ResolveURL(stream)
		check(err)

		filename := normalizeFileName(player.Title()) + "." + stream.FileExtension()

		fmt.Printf("Stream URL: %q\n\nDownloading %q...\n", url, filename)

		check(nicehttp.DownloadFile(filename, url))
	}
}
