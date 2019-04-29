package poller

import (
	"fmt"
	"github.com/larryrun/tor-poller/pkg/config"
	"github.com/larryrun/tor-poller/pkg/download"
	"github.com/larryrun/tor-poller/pkg/mags"
	"log"
	"os"
	"path"
	"time"
)

type Poller struct {
	Conf *config.Config
}

func NewPoller(conf *config.Config) *Poller {
	return &Poller{Conf: conf}
}

func (p *Poller) StartPolling() {
	func() {
		p.pollOnce()
		pollingTimer := time.NewTimer(2 * time.Hour)
		for {
			<-pollingTimer.C
			p.pollOnce()
		}
	}()
}

func (p *Poller) pollOnce() {
	items := PrepareAvailableItems(p.Conf)
	for _, item := range items {
		err := download.NewDownload(p.Conf, item, path.Join(p.Conf.DownloadFolder, item.Name))
		if err != nil {
			if err.Error() != download.ErrorSameEpisodeDownloadExists {
				log.Println(err.Error())
			}
		}
	}
}

func PrepareAvailableItems(conf *config.Config) []*mags.TargetItem {
	finders := make([]mags.MagFinder, 0)
	for _, e := range conf.Items {
		destFolder, err := EnsureDestFolder(&e, conf)
		if err != nil {
			log.Println(err)
			continue
		}
		if e.Type == "itvfans" {
			itvFansFinder := &mags.ITVFansFinder{Name: e.Name, Page: e.Info, DownloadFolder: destFolder, Extra: e.Extra}
			finders = append(finders, itvFansFinder)
		}
	}

	totalItems := make([]*mags.TargetItem, 0)
	for _, finder := range finders {
		targetItems, err := finder.ListAvailableItems()
		if err != nil {
			log.Println(err.Error())
		} else {
			for _, targetItem := range targetItems {
				totalItems = append(totalItems, targetItem)
			}
		}
	}
	return totalItems
}

func EnsureDestFolder(item *config.ItemInfo, config *config.Config) (string, error) {
	itemDownloadFolder := path.Join(config.DownloadFolder, item.Name)
	if _, err := os.Stat(itemDownloadFolder); os.IsNotExist(err) {
		err := os.MkdirAll(itemDownloadFolder, 0777)
		if err != nil {
			return "", fmt.Errorf("failed to create download folder: %s, cause: %s", itemDownloadFolder, err.Error())
		}
	}
	return itemDownloadFolder, nil
}
