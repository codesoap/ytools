package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools"
	"os"
)

const max_results = 16

type Video struct {
	Title string
	Url   string
}

type YtInitialData struct {
	Contents struct {
		TwoColumnWatchNextResults struct {
			SecondaryResults struct {
				SecondaryResults struct {
					Results []struct {
						CompactVideoRenderer CompactVideoRenderer
					}
				}
			}
		}
	}
}

type CompactVideoRenderer struct {
	VideoId string
	Title   struct {
		SimpleText string
	}
}

func main() {
	video_url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the video URL: %s\n", err.Error())
		os.Exit(1)
	}
	recommendations, err := scrape_off_recommendations(video_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find recommendations: %s\n", err.Error())
		os.Exit(1)
	}
	if len(recommendations) == 0 {
		fmt.Fprintf(os.Stderr, "No recommendations found.\n")
		os.Exit(1)
	}
	if err := save_recommendations_urls(recommendations); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save found URLs: %s\n", err.Error())
		os.Exit(1)
	}
	print_video_titles(recommendations)
}

func scrape_off_recommendations(video_url string) (videos []Video, err error) {
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(video_url); err != nil {
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

	sr := data.Contents.TwoColumnWatchNextResults.SecondaryResults
	for _, result := range sr.SecondaryResults.Results {
		var video Video
		video, err = extract_video_from_video_renderer(result.CompactVideoRenderer)
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
	return
}

func extract_video_from_video_renderer(renderer CompactVideoRenderer) (video Video, err error) {
	if len(renderer.VideoId) == 0 {
		err = fmt.Errorf("videoId is missing in videoRenderer")
		return
	}
	if len(renderer.Title.SimpleText) == 0 {
		err = fmt.Errorf("no title found for videoRenderer")
		return
	}
	video = Video{
		Url:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", renderer.VideoId),
		Title: renderer.Title.SimpleText,
	}
	return
}

func save_recommendations_urls(videos []Video) (err error) {
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
