package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codesoap/ytools"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const max_results = 12

type Video struct {
	Title string
	Url   string
}

// TODO: Consider using nested structs to save some lines of code.

type YtInitialData struct {
	Contents YtInitialDataContents
}

type YtInitialDataContents struct {
	TwoColumnSearchResultsRenderer TwoColumnSearchResultsRenderer
}

type TwoColumnSearchResultsRenderer struct {
	PrimaryContents PrimaryContents
}

type PrimaryContents struct {
	SectionListRenderer SectionListRenderer
}

type SectionListRenderer struct {
	Contents []SectionListRendererContent
}

type SectionListRendererContent struct {
	ItemSectionRenderer ItemSectionRenderer
}

type ItemSectionRenderer struct {
	Contents []ItemSectionRendererContent
}

type ItemSectionRendererContent struct {
	// ShelfRenderer ShelfRenderer
	VideoRenderer VideoRenderer
}

// type ShelfRenderer struct {
// 	Content ShelfRendererContent
// }

// type ShelfRendererContent struct {
// 	VerticalListRenderer VerticalListRenderer
// }

// type VerticalListRenderer struct {
// 	Items []VerticalListRendererItem
// }

// type VerticalListRendererItem struct {
// 	VideoRenderer VideoRenderer
// }

type VideoRenderer struct {
	VideoId string
	Title   VideoRendererTitle
}

type VideoRendererTitle struct {
	Runs []Run
}

type Run struct {
	Text string
}

func main() {
	search_url := get_search_url()
	videos, err := scrape_off_videos(search_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	if len(videos) == 0 {
		fmt.Fprintf(os.Stderr, "No videos found.\n")
		os.Exit(1)
	}
	if err := save_videos_urls(videos); err != nil {
		fmt.Fprintf(os.Stderr, "Failed saving found URLs.\n")
		os.Exit(1)
	}
	print_video_titles(videos)
}

func get_search_url() string {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Give one or more search terms as parameters.\n")
		os.Exit(1)
	}
	search_string := url.PathEscape(strings.Join(os.Args[1:], " "))
	return fmt.Sprintf(
		"https://www.youtube.com/results?search_query=%s",
		search_string)
}

func scrape_off_videos(search_url string) (videos []Video, err error) {
	resp, err := http.Get(search_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get '%s'\n", search_url)
		return
	}
	defer resp.Body.Close()
	dataJson, err := extract_json(resp.Body)
	if err != nil {
		return
	}
	return extract_videos(dataJson)
}

func extract_json(body io.Reader) (mainJson []byte, err error) {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`(?m)^ *window\["ytInitialData"\] *= *(.*); *$`)
	matches := re.FindSubmatch(bytes)
	if matches == nil {
		err = errors.New("retrieved HTML does not contain the expected JSON")
		return
	}
	return matches[1], nil
}

func extract_videos(dataJson []byte) (videos []Video, err error) {
	videos = make([]Video, 0, max_results)
	var data YtInitialData
	if err = json.Unmarshal(dataJson, &data); err != nil {
		return
	}

	pc := data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents
	for _, slrContent := range pc.SectionListRenderer.Contents {
		for _, isrContent := range slrContent.ItemSectionRenderer.Contents {
			var video Video
			video, err = extract_video_from_video_renderer(isrContent.VideoRenderer)
			if err != nil {
				// This sometimes happens, but I don't think it's problematic.
				err = nil
				continue
			}
			videos = append(videos, video)
			if len(videos) == max_results {
				break
			}
		}
	}
	return
}

func extract_video_from_video_renderer(renderer VideoRenderer) (video Video, err error) {
	if len(renderer.VideoId) == 0 {
		err = errors.New("videoId is missing in videoRenderer")
		return
	}
	if len(renderer.Title.Runs) != 1 {
		err = errors.New("multiple or no runs found for videoRenderer")
		return
	}
	if len(renderer.Title.Runs[0].Text) == 0 {
		err = errors.New("no title found for videoRenderer")
		return
	}
	video = Video{
		Url: fmt.Sprintf("https://www.youtube.com/watch?v=%s", renderer.VideoId),
		Title: renderer.Title.Runs[0].Text,
	}
	return
}

func save_videos_urls(videos []Video) (err error) {
	videos_urls := make([]string, 0, max_results)
	for _, video := range videos {
		videos_urls = append(videos_urls, video.Url)
	}
	return ytools.SaveUrls(videos_urls)
}

func print_video_titles(videos []Video) {
	for i, video := range videos {
		fmt.Printf("%2d: %s\n", i+1, video.Title)
	}
}
