package main

import (
	"encoding/json"
	"fmt"
	"github.com/codesoap/ytools"
	"os"
)

type Info struct {
	Title string
	Views string
	// TODO: Length string
	// The length is not available on the main page, will probably have
	// to load something like this:
	// https://www.youtube.com/annotations_invideo?video_id=DuoTdnq_OqE
	Likes       string
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
	VideoActions struct {
		MenuRenderer struct {
			TopLevelButtons []struct {
				SegmentedLikeDislikeButtonViewModel struct {
					LikeButtonViewModel struct {
						LikeButtonViewModel struct {
							ToggleButtonViewModel struct {
								ToggleButtonViewModel struct {
									DefaultButtonViewModel struct {
										ButtonViewModel struct {
											Title string
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	DateText struct {
		SimpleText string
	}
}

type VideoSecondaryInfoRenderer struct {
	AttributedDescription struct {
		Content string
	}
}

func main() {
	videoUrl, err := ytools.GetDesiredVideoUrl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get video URL: %s\n", err.Error())
		os.Exit(1)
	}
	info, err := scrapeOffInfo(videoUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to scrape the videos page:", err.Error())
		os.Exit(1)
	}
	printInfo(info)
}

func printInfo(info Info) {
	fmt.Println(info.Title)
	if info.Description == "" {
		f := "%s  %s  %s\n"
		fmt.Printf(f, info.Views, info.Likes, info.Date)
	} else {
		f := "%s  %s  %s\n\n"
		fmt.Printf(f, info.Views, info.Likes, info.Date)
		fmt.Println(info.Description)
	}
}

func scrapeOffInfo(url string) (info Info, err error) {
	var dataJson []byte
	if dataJson, err = ytools.ExtractJson(url); err != nil {
		return
	}
	return extractInfo(dataJson)
}

func extractInfo(dataJson []byte) (info Info, err error) {
	var data YtInitialData
	if err = json.Unmarshal(dataJson, &data); err != nil {
		return
	}
	r := data.Contents.TwoColumnWatchNextResults.Results.Results
	if len(r.Contents) == 0 {
		return info, fmt.Errorf("no contents found in JSON")
	}
	primaryInfo := r.Contents[0].VideoPrimaryInfoRenderer
	secondaryInfo := r.Contents[1].VideoSecondaryInfoRenderer
	if err = fillTitle(&info, primaryInfo); err != nil {
		return
	}
	if err = fillViews(&info, primaryInfo); err != nil {
		return
	}
	if err = fillLikes(&info, primaryInfo); err != nil {
		return
	}
	if err = fillDate(&info, primaryInfo); err != nil {
		return
	}
	// TODO: Owner
	fillDescription(&info, secondaryInfo)
	return
}

func fillTitle(info *Info, data VideoPrimaryInfoRenderer) error {
	if len(data.Title.Runs) == 0 {
		return fmt.Errorf("no found in JSON")
	}
	// There are multiple runs when the title contains hashtags.
	// Just join the parts together.
	info.Title = ""
	for _, r := range data.Title.Runs {
		info.Title += r.Text
	}
	if len(info.Title) == 0 {
		return fmt.Errorf("title is empty")
	}
	return nil
}

func fillViews(info *Info, data VideoPrimaryInfoRenderer) error {
	info.Views = data.ViewCount.VideoViewCountRenderer.ViewCount.SimpleText
	if len(info.Views) == 0 {
		return fmt.Errorf("views is empty")
	}
	return nil
}

func fillLikes(info *Info, data VideoPrimaryInfoRenderer) error {
	for _, button := range data.VideoActions.MenuRenderer.TopLevelButtons {
		info.Likes = button.
			SegmentedLikeDislikeButtonViewModel.
			LikeButtonViewModel.
			LikeButtonViewModel.
			ToggleButtonViewModel.
			ToggleButtonViewModel.
			DefaultButtonViewModel.
			ButtonViewModel.
			Title
		if info.Likes != "" {
			info.Likes += " ▲"
			return nil
		}
	}
	return fmt.Errorf("like count not found")
}

func fillDate(info *Info, data VideoPrimaryInfoRenderer) error {
	info.Date = data.DateText.SimpleText
	if len(info.Date) == 0 {
		return fmt.Errorf("date is empty")
	}
	return nil
}

func fillDescription(info *Info, data VideoSecondaryInfoRenderer) {
	info.Description = data.AttributedDescription.Content
}
