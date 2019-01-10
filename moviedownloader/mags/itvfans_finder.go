package mags

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ITVFansFinder struct {
	Name           string
	Page           string
	DownloadFolder string
}

var titlePattern *regexp.Regexp
func init() {
	titlePattern, _ = regexp.Compile("\\.S(\\d+)E(\\d+)\\.")
}

func (finder *ITVFansFinder) ListAvailableItems() ([]*TargetItem, error) {
	collector := colly.NewCollector(colly.Async(true))
	collector.SetRequestTimeout(time.Duration(time.Second * 10))

	itemChan := make(chan *TargetItem)
	defer close(itemChan)

	items := make([]*TargetItem, 1)
	mutex := sync.Mutex{}
	collector.OnHTML("div.wpb_wrapper li", func(e *colly.HTMLElement) {
		item := &TargetItem{}
		title := strings.TrimSpace(e.Attr("title"))
		if title != "" {
			magLinkDom := e.DOM.Find("a[href^=magnet]")
			link, exists := magLinkDom.Attr("href")
			if exists {
				item.Link = link
			} else {
				return
			}

			item.FileName = title
			err := itvfans_fill_episode_resolution_info(item, title)
			if err != nil {
				log.Printf("Failed to get episode info from: %s", title)
				return
			}

			sizeAndUnitDom := e.DOM.ChildrenFiltered("font")
			if sizeAndUnitDom != nil {
				sizeAndUnit := sizeAndUnitDom.Text()
				if strings.TrimSpace(sizeAndUnit) != "" {
					err = itvfans_fill_size_info(item, sizeAndUnit[1:len(sizeAndUnit)-1])
					if err != nil {
						log.Printf("Failed to get size info from: %s", sizeAndUnit)
					}
				}
			}

			mutex.Lock()
			log.Println(item.Link)
			items = append(items, item)
			mutex.Unlock()
		}
	})
	if err := collector.Visit(finder.Page); err != nil {
		return nil, fmt.Errorf("failed to list iTVFans for %s, cause: %s", finder.Name, err.Error())
	}
	collector.Wait()
	return items, nil
}

func itvfans_fill_size_info(item *TargetItem, sizeAndUnit string) error {
	var err error
	if strings.HasSuffix(sizeAndUnit, "GB") {
		sizeStr := strings.Replace(sizeAndUnit, "GB", "", -1)
		size, err := strconv.ParseFloat(sizeStr, 32)
		if err == nil {
			item.Size = int(size * 1024)
		}
	} else if strings.HasSuffix(sizeAndUnit, "MB") {
		sizeStr := strings.Replace(sizeAndUnit, "MB", "", -1)
		size, err := strconv.ParseFloat(sizeStr, 32)
		if err == nil {
			item.Size = int(size)
		}
	}
	if item.Size == 0 {
		return fmt.Errorf("failed to parse Size info from: %s", sizeAndUnit)
	} else if err != nil {
		return fmt.Errorf("failed to parse Size info from: %s, cause: %s", sizeAndUnit, err.Error())
	}
	return nil
}

func itvfans_fill_episode_resolution_info(item *TargetItem, title string) error {
	upperTitle := strings.ToUpper(title)
	if strings.Contains(upperTitle, ".1080P.") {
		item.Resolution = "1080p"
	}
	if strings.Contains(upperTitle, ".720P.") {
		item.Resolution = "720p"
	}
	submatch := titlePattern.FindStringSubmatch(upperTitle)
	if len(submatch) == 3 {
		session, err := strconv.ParseInt(submatch[1], 10, 32)
		if err == nil {
			episode, err := strconv.ParseInt(submatch[2], 10, 32)
			if err == nil {
				item.Session = int(session)
				item.Episode = int(episode)
			} else {
				return fmt.Errorf("failed to parse Episode info from title: %s, cause: %s", title, err.Error())
			}
		} else {
			return fmt.Errorf("failed to parse Session info from title: %s, cause: %s", title, err.Error())
		}
	}
	return nil
}