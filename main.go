package main

import (
	"github.com/larryrun/movie-downloader/moviedownloader"
	"github.com/larryrun/movie-downloader/moviedownloader/mags"
	"log"
)

func main() {
	moviedownloader.SetConfigPath("./cfg.yaml")
	PrepareAvailableItems()
}

func PrepareAvailableItems() {
	config, err := moviedownloader.ReadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	finders := make([]mags.MagFinder, len(config.Items))
	for i, e := range config.Items {
		if e.Type == "itvfans" {
			itvFansFinder := &mags.ITVFansFinder{Name: e.Name, Page: e.Info, DownloadFolder: ""}
			finders[i] = itvFansFinder
		}

		for _, finder := range finders {
			targetItems, err := finder.ListAvailableItems()
			if err != nil {
				log.Println(err.Error())
			} else {
				for _, targetItem := range targetItems {
					log.Printf("%+v", targetItem)
				}
			}
		}
	}
}