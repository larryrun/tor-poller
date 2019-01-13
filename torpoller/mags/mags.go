package mags

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TargetItem struct {
	Name       string
	Episode    Episode
	Size       int
	Link       string
	Resolution string
}

type Episode struct {
	Season  int
	Episode int
}

type MagFinder interface {
	ListAvailableItems() ([]*TargetItem, error)
}

func (thisEpisode *Episode) String() string {
	return fmt.Sprintf("%d-%d", thisEpisode.Season, thisEpisode.Episode)
}

func (thisEpisode *Episode) laterThan(anotherEpisode *Episode) bool {
	if thisEpisode.Season > anotherEpisode.Season {
		return true
	} else if thisEpisode.Season < anotherEpisode.Season {
		return false
	} else {
		return thisEpisode.Episode > anotherEpisode.Episode
	}
}

func (thisEpisode *Episode) earlierThan(anotherEpisode *Episode) bool {
	if thisEpisode.Season == anotherEpisode.Season && thisEpisode.Episode == anotherEpisode.Episode {
		return false
	} else {
		return !thisEpisode.laterThan(anotherEpisode)
	}
}

var episodePattern *regexp.Regexp
func init() {
	episodePattern, _ = regexp.Compile("\\WS(\\d+)E(\\d+)\\W")
}

func ExtractEpisodeInfo(title string)(Episode, error) {
	upperTitle := strings.ToUpper(title)
	subMatch := episodePattern.FindStringSubmatch(upperTitle)
	if len(subMatch) == 3 {
		season, err := strconv.Atoi(subMatch[1])
		if err == nil {
			episode, err := strconv.Atoi(subMatch[2])
			if err == nil {
				if episode > 100 {
					return Episode{}, fmt.Errorf("the extracted Episode is: %d, which is very likely a false result, from title: %s", episode, title)
				} else {
					return Episode{Season: season, Episode: episode}, nil
				}
			}
		}
	}
	return Episode{}, fmt.Errorf("failed to extract Episode info from title: %s", title)
}

func ExtractResolutionInfo(title string) (string, error) {
	upperTitle := strings.ToUpper(title)
	if strings.Contains(upperTitle, ".1080P.") {
		return "1080p", nil
	}
	if strings.Contains(upperTitle, ".720P.") {
		return "720p", nil
	}
	return "", fmt.Errorf("cannot extract resolution info from title: %s", title)
}

func ExtractItemSize(sizeAndUnit string) (int, error) {
	err := fmt.Errorf("")
	if strings.HasSuffix(sizeAndUnit, "GB") {
		sizeStr := strings.Replace(sizeAndUnit, "GB", "", -1)
		size, err := strconv.ParseFloat(sizeStr, 32)
		if err == nil {
			return int(size * 1024), nil
		}
	} else if strings.HasSuffix(sizeAndUnit, "MB") {
		sizeStr := strings.Replace(sizeAndUnit, "MB", "", -1)
		size, err := strconv.ParseFloat(sizeStr, 32)
		if err == nil {
			return int(size), nil
		}
	}
	return 0, fmt.Errorf("failed to parse Size info from: %s, cause: %s", sizeAndUnit, err.Error())
}

func GetLatestEpisodeInTitles(titles []string) Episode {
	var latestEpisode *Episode
	for _, title := range titles {
		episode, err := ExtractEpisodeInfo(title)
		if err == nil {
			if latestEpisode == nil || episode.laterThan(latestEpisode) {
				latestEpisode = &episode
			}
		}
	}
	if latestEpisode == nil {
		return Episode{}
	} else {
		return *latestEpisode
	}
}
