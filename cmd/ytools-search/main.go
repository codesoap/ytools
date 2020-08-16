package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools"
	"net/url"
	"os"
	"strings"
)

const max_results = 12

type Video struct {
	Title string
	Url   string
}

type YtInitialData struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								VideoRenderer VideoRenderer
							}
						}
					}
				}
			}
		}
	}
}

type VideoRenderer struct {
	VideoId string
	Title   struct {
		Runs []struct {
			Text string
		}
	}
}

func main() {
	search_url := get_search_url()
	videos, err := scrape_off_videos(search_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get video URLs: %s\n", err.Error())
		os.Exit(1)
	}
	if len(videos) == 0 {
		fmt.Fprintf(os.Stderr, "No videos found.\n")
		os.Exit(1)
	}
	if err := save_videos_urls(videos); err != nil {
		fmt.Fprintf(os.Stderr, "Failed saving found URLs: %s\n", err.Error())
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
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(search_url); err != nil {
		return
	}
	return extract_videos(dataJson)
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
		err = fmt.Errorf("videoId is missing in videoRenderer")
		return
	}
	if len(renderer.Title.Runs) != 1 {
		err = fmt.Errorf("multiple or no runs found for videoRenderer")
		return
	}
	if len(renderer.Title.Runs[0].Text) == 0 {
		err = fmt.Errorf("no title found for videoRenderer")
		return
	}
	video = Video{
		Url:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", renderer.VideoId),
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
