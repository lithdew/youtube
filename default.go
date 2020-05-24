// Package youtube is a library for retrieving metadata and obtaining direct links to video-only/audio-only/muxed
// streams of videos on YouTube.

package youtube

import "time"

var defaultClient = NewClient()

func Load(id StreamID) (Player, error) {
	return defaultClient.Load(id)
}

func LoadTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return defaultClient.LoadTimeout(id, timeout)
}

func LoadDeadline(id StreamID, deadline time.Time) (Player, error) {
	return defaultClient.LoadDeadline(id, deadline)
}

func LoadPlaylist(id string, offset uint) (PlaylistResult, error) {
	return defaultClient.LoadPlaylist(id, offset)
}

func LoadPlaylistTimeout(id string, offset uint, timeout time.Duration) (PlaylistResult, error) {
	return defaultClient.LoadPlaylistTimeout(id, offset, timeout)
}

func LoadPlaylistDeadline(id string, offset uint, deadline time.Time) (PlaylistResult, error) {
	return defaultClient.LoadPlaylistDeadline(id, offset, deadline)
}

func Search(query string, page uint) (SearchResult, error) {
	return defaultClient.Search(query, page)
}

func SearchTimeout(query string, page uint, timeout time.Duration) (SearchResult, error) {
	return defaultClient.SearchTimeout(query, page, timeout)
}

func SearchDeadline(query string, page uint, deadline time.Time) (SearchResult, error) {
	return defaultClient.SearchDeadline(query, page, deadline)
}

func LoadWatchPlayer(id StreamID) (Player, error) {
	return defaultClient.LoadWatchPlayer(id)
}

func LoadWatchPlayerTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return defaultClient.LoadWatchPlayerTimeout(id, timeout)
}

func LoadWatchPlayerDeadline(id StreamID, deadline time.Time) (Player, error) {
	return defaultClient.LoadWatchPlayerDeadline(id, deadline)
}

func LoadEmbedPlayerAssets(id StreamID) (Assets, error) {
	return defaultClient.LoadEmbedPlayerAssets(id)
}

func LoadEmbedPlayerAssetsTimeout(id StreamID, timeout time.Duration) (Assets, error) {
	return defaultClient.LoadEmbedPlayerAssetsTimeout(id, timeout)
}

func LoadEmbedPlayerAssetsDeadline(id StreamID, deadline time.Time) (Assets, error) {
	return defaultClient.LoadEmbedPlayerAssetsDeadline(id, deadline)
}

func LoadEmbedPlayerStreams(id StreamID) (Streams, error) {
	return defaultClient.LoadEmbedPlayerStreams(id)
}

func LoadEmbedPlayerStreamsTimeout(id StreamID, timeout time.Duration) (Streams, error) {
	return defaultClient.LoadEmbedPlayerStreamsTimeout(id, timeout)
}

func LoadEmbedPlayerStreamsDeadline(id StreamID, deadline time.Time) (Streams, error) {
	return defaultClient.LoadEmbedPlayerStreamsDeadline(id, deadline)
}

func LoadEmbedPlayer(id StreamID) (Player, error) {
	return defaultClient.LoadEmbedPlayer(id)
}

func LoadEmbedPlayerTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return defaultClient.LoadEmbedPlayerTimeout(id, timeout)
}

func LoadEmbedPlayerDeadline(id StreamID, deadline time.Time) (Player, error) {
	return defaultClient.LoadEmbedPlayerDeadline(id, deadline)
}
