package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools"
	"net/url"
	"os"
	"strings"
)

const maxResults = 12

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
	searchUrl := getSearchUrl()
	videos, err := scrapeOffVideos(searchUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get video URLs: %s\n", err.Error())
		os.Exit(1)
	}
	if len(videos) == 0 {
		fmt.Fprintf(os.Stderr, "No videos found.\n")
		os.Exit(1)
	}
	if err := saveVideosUrls(videos); err != nil {
		fmt.Fprintf(os.Stderr, "Failed saving found URLs: %s\n", err.Error())
		os.Exit(1)
	}
	printVideoTitles(videos)
}

func getSearchUrl() string {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Give one or more search terms as parameters.\n")
		os.Exit(1)
	}
	searchString := url.QueryEscape(strings.Join(os.Args[1:], " "))
	return fmt.Sprintf(
		"https://www.youtube.com/results?search_query=%s",
		searchString)
}

func scrapeOffVideos(searchUrl string) (videos []Video, err error) {
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(searchUrl); err != nil {
		return
	}
	return extractVideos(dataJson)
}

func extractVideos(dataJson []byte) (videos []Video, err error) {
	videos = make([]Video, 0, maxResults)
	var data YtInitialData
	if err = json.Unmarshal(dataJson, &data); err != nil {
		return
	}

	pc := data.Contents.TwoColumnSearchResultsRenderer.PrimaryContents
	for _, slrContent := range pc.SectionListRenderer.Contents {
		for _, isrContent := range slrContent.ItemSectionRenderer.Contents {
			var video Video
			video, err = extractVideoFromVideoRenderer(isrContent.VideoRenderer)
			if err != nil {
				// This sometimes happens, but I don't think it's problematic.
				err = nil
				continue
			}
			videos = append(videos, video)
			if len(videos) == maxResults {
				break
			}
		}
	}
	return
}

func extractVideoFromVideoRenderer(renderer VideoRenderer) (video Video, err error) {
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

func saveVideosUrls(videos []Video) (err error) {
	videosUrls := make([]string, 0, maxResults)
	for _, video := range videos {
		videosUrls = append(videosUrls, video.Url)
	}
	return ytools.SaveUrls(videosUrls)
}

func printVideoTitles(videos []Video) {
	for i, video := range videos {
		fmt.Printf("%2d: %s\n", i+1, video.Title)
	}
}
