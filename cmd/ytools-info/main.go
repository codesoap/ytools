package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codesoap/ytools"
	"os"
	"strings"
)

type Info struct {
	Title string
	Views string
	// TODO: Length string
	// The length is not available on the main page, will probably have
	// to load something like this:
	// https://www.youtube.com/annotations_invideo?video_id=DuoTdnq_OqE
	Likes       string
	Dislikes    string
	Date        string
	Description string
}

type YtInitialData struct {
	Contents struct {
		TwoColumnWatchNextResults struct {
			Results struct {
				Results struct {
					Contents []struct {
						VideoPrimaryInfoRenderer   VideoPrimaryInfoRenderer
						VideoSecondaryInfoRenderer VideoSecondaryInfoRenderer
					}
				}
			}
		}
	}
}

type VideoPrimaryInfoRenderer struct {
	Title struct {
		Runs []struct {
			Text string
		}
	}
	ViewCount struct {
		VideoViewCountRenderer struct {
			ViewCount struct {
				SimpleText string
			}
		}
	}
	SentimentBar struct {
		SentimentBarRenderer struct {
			Tooltip string
		}
	}
	DateText struct {
		SimpleText string
	}
}

type VideoSecondaryInfoRenderer struct {
	Owner struct {
		VideoOwnerRenderer struct {
			Title struct {
				Runs []struct {
					Text string
				}
			}
		}
	}
	Description struct {
		Runs []struct {
			Text string
		}
	}
}

func main() {
	video_url, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		os.Exit(1)
	}
	info, err := scrape_off_info(video_url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to scrape the videos page")
		os.Exit(1)
	}
	print_info(info)
}

func print_info(info Info) {
	fmt.Println(info.Title)
	f := "%s  ▲ %s  ▼ %s  %s\n\n"
	fmt.Printf(f, info.Views, info.Likes, info.Dislikes, info.Date)
	fmt.Println(info.Description)
}

func scrape_off_info(url string) (info Info, err error) {
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(url); err != nil {
		return
	}
	return extract_info(dataJson)
}

func extract_info(dataJson []byte) (info Info, err error) {
	var data YtInitialData
	if err = json.Unmarshal(dataJson, &data); err != nil {
		return
	}
	r := data.Contents.TwoColumnWatchNextResults.Results.Results
	if len(r.Contents) == 0 {
		return info, errors.New("no contents found in JSON")
	}
	primaryInfo := r.Contents[0].VideoPrimaryInfoRenderer
	secondaryInfo := r.Contents[1].VideoSecondaryInfoRenderer
	if err = fill_title(&info, primaryInfo); err != nil {
		return
	}
	if err = fill_views(&info, primaryInfo); err != nil {
		return
	}
	if err = fill_votes(&info, primaryInfo); err != nil {
		return
	}
	if err = fill_date(&info, primaryInfo); err != nil {
		return
	}
	// TODO: Owner
	fill_description(&info, secondaryInfo)
	return
}

func fill_title(info *Info, data VideoPrimaryInfoRenderer) error {
	if len(data.Title.Runs) != 1 {
		return errors.New("no or multiple titles found in JSON")
	}
	info.Title = data.Title.Runs[0].Text
	if len(info.Title) == 0 {
		return errors.New("title is empty")
	}
	return nil
}

func fill_views(info *Info, data VideoPrimaryInfoRenderer) error {
	info.Views = data.ViewCount.VideoViewCountRenderer.ViewCount.SimpleText
	if len(info.Views) == 0 {
		return errors.New("views is empty")
	}
	return nil
}

func fill_votes(info *Info, data VideoPrimaryInfoRenderer) error {
	tooltip := data.SentimentBar.SentimentBarRenderer.Tooltip
	numbers := strings.Split(tooltip, " / ")
	if len(numbers) != 2 {
		return errors.New("could not parse likes and dislikes")
	}
	info.Likes = numbers[0]
	info.Dislikes = numbers[1]
	return nil
}

func fill_date(info *Info, data VideoPrimaryInfoRenderer) error {
	info.Date = data.DateText.SimpleText
	if len(info.Date) == 0 {
		return errors.New("date is empty")
	}
	return nil
}

func fill_description(info *Info, data VideoSecondaryInfoRenderer) {
	desc := []byte("")
	for _, run := range data.Description.Runs {
		desc = append(desc, []byte(run.Text)...)
	}
	info.Description = string(desc)
}
