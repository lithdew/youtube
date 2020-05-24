package youtube

import (
	"errors"
	"fmt"
	"github.com/lithdew/bytesutil"
	"github.com/lithdew/nicehttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"golang.org/x/sync/errgroup"
	"net/url"
	"time"
)

var zeroTime time.Time

// Transport represents the transport client used for youtube.Client.
type Transport interface {
	DownloadBytesDeadline(dst []byte, url string, deadline time.Time) ([]byte, error)
}

type Client struct {
	Transport
}

func NewClient() Client {
	client := nicehttp.NewClient()
	return WrapClient(&client)
}

func WrapClient(transport Transport) Client {
	return Client{Transport: transport}
}

func (c *Client) Load(id StreamID) (Player, error) {
	return c.LoadDeadline(id, zeroTime)
}

func (c *Client) LoadTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return c.LoadDeadline(id, time.Now().Add(timeout))
}

func (c *Client) LoadDeadline(id StreamID, deadline time.Time) (Player, error) {
	var (
		player Player
		err    error
	)

	if err = id.Valid(); err != nil {
		return player, err
	}

	// Attempt to grab the standard player first.

	player, err = c.LoadWatchPlayerDeadline(id, deadline)
	if err == nil {
		return player, nil
	}

	// If it fails, attempt to grab the embedded player second.

	player, err = c.LoadEmbedPlayerDeadline(id, deadline)
	if err == nil {
		return player, nil
	}

	// If all fails, throw an error :c.

	return player, fmt.Errorf("failed to load player: %w", err)
}

func (c *Client) LoadPlaylist(id string, offset uint) (PlaylistResult, error) {
	return c.LoadPlaylistDeadline(id, offset, zeroTime)
}

func (c *Client) LoadPlaylistTimeout(id string, offset uint, timeout time.Duration) (PlaylistResult, error) {
	return c.LoadPlaylistDeadline(id, offset, time.Now().Add(timeout))
}

func (c *Client) LoadPlaylistDeadline(id string, offset uint, deadline time.Time) (PlaylistResult, error) {
	var result PlaylistResult

	uri := []byte("https://www.youtube.com/list_ajax?style=json&action_get_list=1")

	uri = append(uri, "&list="...)
	uri = append(uri, id...)

	uri = append(uri, "&index="...)
	uri = fasthttp.AppendUint(uri, int(offset))

	uri = append(uri, "&hl="...)
	uri = append(uri, "en"...)

	buf, err := c.DownloadBytesDeadline(nil, bytesutil.String(uri), deadline)
	if err != nil {
		return result, fmt.Errorf("failed to load offset %d of playlist %q: %w", offset, id, err)
	}

	val, err := fastjson.ParseBytes(buf)
	if err != nil {
		return result, fmt.Errorf("got malformed json loading offset %d of playlist %q: %w", offset, id, err)
	}

	return ParsePlaylistResultJSON(val), nil
}

func (c *Client) Search(query string, page uint) (SearchResult, error) {
	return c.SearchDeadline(query, page, zeroTime)
}

func (c *Client) SearchTimeout(query string, page uint, timeout time.Duration) (SearchResult, error) {
	return c.SearchDeadline(query, page, time.Now().Add(timeout))
}

func (c *Client) SearchDeadline(query string, page uint, deadline time.Time) (SearchResult, error) {
	var result SearchResult

	uri := []byte("https://www.youtube.com/search_ajax?style=json")

	uri = append(uri, "&search_query="...)
	uri = append(uri, url.PathEscape(query)...)

	uri = append(uri, "&page="...)
	uri = fasthttp.AppendUint(uri, int(page))

	uri = append(uri, "&hl="...)
	uri = append(uri, "en"...)

	buf, err := c.DownloadBytesDeadline(nil, bytesutil.String(uri), deadline)
	if err != nil {
		return result, fmt.Errorf("failed to search for page %d of query %q: %w", page, query, err)
	}

	val, err := fastjson.ParseBytes(buf)
	if err != nil {
		return result, fmt.Errorf("got malformed json searching for page %d of query %q: %w", page, query, err)
	}

	return ParseSearchResultJSON(val), nil
}

func (c *Client) LoadWatchPlayer(id StreamID) (Player, error) {
	return c.LoadWatchPlayerDeadline(id, zeroTime)
}

func (c *Client) LoadWatchPlayerTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return c.LoadWatchPlayerDeadline(id, time.Now().Add(timeout))
}

func (c *Client) LoadWatchPlayerDeadline(id StreamID, deadline time.Time) (Player, error) {
	player := Player{Transport: c.Transport}

	if err := id.Valid(); err != nil {
		return player, err
	}

	// Download player HTML.

	buf, err := c.DownloadBytesDeadline(nil, "https://www.youtube.com/watch?v="+string(id), deadline)
	if err != nil {
		return player, err
	}

	matches := RegexWatchPlayerConfig.FindSubmatch(buf)
	if matches == nil {
		return player, errors.New("could not find watch video player config in html page")
	}

	val, err := fastjson.ParseBytes(matches[1])
	if err != nil {
		return player, fmt.Errorf("failed to parse video player config: %w", err)
	}

	player.Assets = ParseAssetsJSON(val.Get("assets"))

	// Extract streaming info.

	val, err = fastjson.ParseBytes(val.GetStringBytes("args", "player_response"))
	if err != nil {
		return player, fmt.Errorf("failed to parse json response: %w", err)
	}

	player.Streams = Streams{v: val}

	if player.Streams.Status() != "OK" {
		return player, fmt.Errorf("unable to get streaming info for id %q: status is %q (reason: %q)", id, player.Streams.Status(), player.Streams.Reason())
	}

	return player, nil
}

func (c *Client) LoadEmbedPlayerAssets(id StreamID) (Assets, error) {
	return c.LoadEmbedPlayerAssetsDeadline(id, zeroTime)
}

func (c *Client) LoadEmbedPlayerAssetsTimeout(id StreamID, timeout time.Duration) (Assets, error) {
	return c.LoadEmbedPlayerAssetsDeadline(id, time.Now().Add(timeout))
}

func (c *Client) LoadEmbedPlayerAssetsDeadline(id StreamID, deadline time.Time) (Assets, error) {
	var assets Assets

	buf, err := c.DownloadBytesDeadline(nil, "https://www.youtube.com/embed/"+string(id), deadline)
	if err != nil {
		return assets, fmt.Errorf("failed to download html of embed player: %w", err)
	}

	matches := RegexEmbedPlayerConfig.FindSubmatch(buf)
	if matches == nil {
		return assets, errors.New("could not find embed player config in html page")
	}

	val, err := fastjson.ParseBytes(matches[1])
	if err != nil {
		return assets, fmt.Errorf("failed to parse embed player config: %w", err)
	}

	return ParseAssetsJSON(val), nil
}

func (c *Client) LoadEmbedPlayerStreams(id StreamID) (Streams, error) {
	return c.LoadEmbedPlayerStreamsDeadline(id, zeroTime)
}

func (c *Client) LoadEmbedPlayerStreamsTimeout(id StreamID, timeout time.Duration) (Streams, error) {
	return c.LoadEmbedPlayerStreamsDeadline(id, time.Now().Add(timeout))
}

func (c *Client) LoadEmbedPlayerStreamsDeadline(id StreamID, deadline time.Time) (Streams, error) {
	var streams Streams

	buf, err := c.DownloadBytesDeadline(nil, "https://www.youtube.com/get_video_info?video_id="+string(id), deadline)
	if err != nil {
		return streams, fmt.Errorf("failed to download stream info: %w", err)
	}

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.ParseBytes(buf)

	status := args.Peek("status")
	if len(status) == 0 {
		return streams, errors.New("key 'status' not found")
	}

	if status := bytesutil.String(status); status != "ok" {
		return streams, fmt.Errorf("status is %q, but expected %q", status, "ok")
	}

	val, err := fastjson.ParseBytes(args.Peek("player_response"))
	if err != nil {
		return streams, fmt.Errorf("failed to parse json response: %w", err)
	}

	streams.v = val

	if streams.Status() != "OK" {
		return streams, fmt.Errorf("unable to get streaming info for id %q: status is %q (reason: %q)", id, streams.Status(), streams.Reason())
	}

	return streams, nil
}

func (c *Client) LoadEmbedPlayer(id StreamID) (Player, error) {
	return c.LoadEmbedPlayerDeadline(id, zeroTime)
}

func (c *Client) LoadEmbedPlayerTimeout(id StreamID, timeout time.Duration) (Player, error) {
	return c.LoadEmbedPlayerDeadline(id, time.Now().Add(timeout))
}

func (c *Client) LoadEmbedPlayerDeadline(id StreamID, deadline time.Time) (Player, error) {
	player := Player{Transport: c.Transport}

	if err := id.Valid(); err != nil {
		return player, err
	}

	var g errgroup.Group

	// Download embed player HTML.

	g.Go(func() error {
		assets, err := c.LoadEmbedPlayerAssetsDeadline(id, deadline)
		if err != nil {
			return err
		}

		player.Assets = assets

		return nil
	})

	// Download streaming info.

	g.Go(func() error {
		streams, err := c.LoadEmbedPlayerStreamsDeadline(id, deadline)
		if err != nil {
			return err
		}

		player.Streams = streams

		return nil
	})

	if err := g.Wait(); err != nil {
		return player, err
	}

	return player, nil
}
