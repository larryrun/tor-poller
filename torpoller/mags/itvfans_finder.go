package mags

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

type ITVFansFinder struct {
	Name           string
	Page           string
	DownloadFolder string
	Extra          map[string]interface{}
}

func (finder *ITVFansFinder) ListAvailableItems() ([]*TargetItem, error) {
	latestEpisode := &Episode{}
	var err error
	v, ok := finder.Extra["start-season"]
	if ok {
		latestEpisode.Season = v.(int)
	} else {
		latestEpisode.Season = 1
	}
	v, ok = finder.Extra["start-episode"]
	if ok {
		latestEpisode.Episode = v.(int) - 1
	} else {
		latestEpisode.Episode = 0
	}
	existingFiles, err := ioutil.ReadDir(finder.DownloadFolder)
	if err != nil {
		log.Printf("failed to read existing files, %s\n", err.Error())
	} else {
		titles := make([]string, len(existingFiles))
		for i := range titles {
			titles[i] = existingFiles[i].Name()
		}
		latestEpisodeInDownloadFolder := GetLatestEpisodeInTitles(titles)
		if latestEpisodeInDownloadFolder.laterThan(latestEpisode) {
			latestEpisode = &latestEpisodeInDownloadFolder
		}
	}

	collector := colly.NewCollector(colly.Async(true)/*, colly.Debugger(&debug.LogDebugger{})*/)
	collector.SetRequestTimeout(time.Duration(time.Second * 10))

	itemChan := make(chan *TargetItem)
	defer close(itemChan)

	items := make([]*TargetItem, 0)
	mutex := sync.Mutex{}
	collector.OnHTML("div.wpb_wrapper li", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.Attr("title"))
		if title != "" {
			episode, err := ExtractEpisodeInfo(title)
			if err != nil {
				log.Println(err)
				return
			}

			if !episode.laterThan(latestEpisode) {
				return
			}

			item := &TargetItem{Episode: episode}
			magLinkDom := e.DOM.Find("a[href^=magnet]")
			link, exists := magLinkDom.Attr("href")
			if exists {
				item.Link = link
			} else {
				return
			}

			resolution, _ := ExtractResolutionInfo(title)
			item.Resolution = resolution

			sizeAndUnitDom := e.DOM.ChildrenFiltered("font")
			if sizeAndUnitDom != nil {
				sizeAndUnit := sizeAndUnitDom.Text()
				if strings.TrimSpace(sizeAndUnit) != "" {
					size, err := ExtractItemSize(sizeAndUnit[1 : len(sizeAndUnit)-1])
					if err != nil {
						log.Println(err)
					}
					item.Size = size
				}
			}
			mutex.Lock()
			items = append(items, item)
			mutex.Unlock()
		}
	})
	if err := collector.Visit(finder.Page); err != nil {
		return nil, fmt.Errorf("failed to list iTVFans for %s, cause: %s", finder.Name, err.Error())
	}
	collector.Wait()

	for i := 1; i < len(items); i++ {
		for j := i; j > 0; j-- {
			if items[j].Episode.earlierThan(&items[j - 1].Episode) {
				temp := items[j]
				items[j] = items[i]
				items[i] = temp
			} else {
				break
			}
		}
	}

	return items, nil
}

