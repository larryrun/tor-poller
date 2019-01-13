package main

import (
	"fmt"
	"github.com/larryrun/tor-poller/torpoller"
	"github.com/larryrun/tor-poller/torpoller/mags"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	torpoller.SetConfigPath("./cfg.yaml")
	DoPolling()
	pollingTimer := time.NewTimer(2 * time.Hour)
	for {
		<-pollingTimer.C
		DoPolling()
	}
}

func DoPolling() {
	config, err := torpoller.ReadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}
	items := PrepareAvailableItems(config)
	for _, item := range items {
		err := torpoller.NewDownload(item, path.Join(config.DownloadFolder, item.Name))
		if err != nil {
			if err.Error() != torpoller.ErrorSameEpisodeDownloadExists {
				log.Println(err.Error())
			}
		}
	}
}

func PrepareAvailableItems(config *torpoller.Config) []*mags.TargetItem {
	finders := make([]mags.MagFinder, 0)
	for _, e := range config.Items {
		destFolder, err := EnsureDestFolder(&e, config)
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

func EnsureDestFolder(item *torpoller.ItemInfo, config *torpoller.Config) (string, error) {
	itemDownloadFolder := path.Join(config.DownloadFolder, item.Name)
	if _, err := os.Stat(itemDownloadFolder); os.IsNotExist(err) {
		err := os.Mkdir(itemDownloadFolder, 0777)
		if err != nil {
			return "", fmt.Errorf("failed to create download folder: %s, cause: %s", itemDownloadFolder, err.Error())
		}
	}
	return itemDownloadFolder, nil
}