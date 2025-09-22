package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools"
	"os"
)

const maxResults = 16

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
						LockupViewModel LockupViewModel
					}
				}
			}
		}
	}
}

type LockupViewModel struct {
	ContentID string
	Metadata  struct {
		LockupMetadataViewModel struct {
			Title struct {
				Content string
			}
		}
	}
}

func main() {
	videoUrl, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the video URL: %s\n", err.Error())
		os.Exit(1)
	}
	recommendations, err := scrapeOffRecommendations(videoUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find recommendations: %s\n", err.Error())
		os.Exit(1)
	}
	if len(recommendations) == 0 {
		fmt.Fprintf(os.Stderr, "No recommendations found.\n")
		os.Exit(1)
	}
	if err := saveRecommendationsUrls(recommendations); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save found URLs: %s\n", err.Error())
		os.Exit(1)
	}
	printVideoTitles(recommendations)
}

func scrapeOffRecommendations(videoUrl string) (videos []Video, err error) {
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(videoUrl); err != nil {
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

	sr := data.Contents.TwoColumnWatchNextResults.SecondaryResults
	for _, result := range sr.SecondaryResults.Results {
		var video Video
		video, err = extractVideoFromVideoModel(result.LockupViewModel)
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
	return
}

func extractVideoFromVideoModel(Model LockupViewModel) (video Video, err error) {
	if len(Model.ContentID) == 0 {
		err = fmt.Errorf("contentId is missing in lockupViewModel")
		return
	}
	if len(Model.Metadata.LockupMetadataViewModel.Title.Content) == 0 {
		err = fmt.Errorf("no title found for lockupViewModel")
		return
	}
	video = Video{
		Url:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", Model.ContentID),
		Title: Model.Metadata.LockupMetadataViewModel.Title.Content,
	}
	return
}

func saveRecommendationsUrls(videos []Video) (err error) {
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
