package torpoller

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/larryrun/tor-poller/torpoller/mags"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

const DownloadTmpFolder = "./download_tmp"
const MaxConcurrentDownloading = 2
const ErrorMaxConcurrentJobReached = "max concurrent job count reached"
const ErrorSameEpisodeDownloadExists = "same episode download exists"

var jobMapMutex sync.Mutex
var downloadingJobMap map[string]*DownloadJob

func init() {
	downloadingJobMap = make(map[string]*DownloadJob, 0)
}

type DownloadJob struct {
	clientCfg  *torrent.ClientConfig
	destFolder string
	item       *mags.TargetItem
}

func (job *DownloadJob) key() string {
	return fmt.Sprintf("%s:%s", job.item.Name, job.item.Episode.String())
}

func NewDownload(item *mags.TargetItem, destFolder string) error {
	jobMapMutex.Lock()
	defer jobMapMutex.Unlock()

	if len(downloadingJobMap) >= MaxConcurrentDownloading {
		return fmt.Errorf(ErrorMaxConcurrentJobReached)
	}
	_, ok := downloadingJobMap[item.Name+item.Episode.String()]
	if ok {
		return fmt.Errorf(ErrorSameEpisodeDownloadExists)
	}
	job, err := createNewJob(item, destFolder)
	if err != nil {
		return err
	}
	go func() {
		err := job.startToDownload()
		delete(downloadingJobMap, job.key())
		if err != nil {
			log.Printf("Failed to download: %s, cause: %s", job.key(), err.Error())
		} else {
			log.Printf("%s has been downloaded successfully", job.key())
		}
	}()
	return nil
}

func createNewJob(item *mags.TargetItem, destFolder string) (*DownloadJob, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = DownloadTmpFolder
	cfg.NoUpload = true
	job := &DownloadJob{destFolder: destFolder, item: item, clientCfg: cfg}
	_, ok := downloadingJobMap[job.key()]
	if ok {
		return nil, fmt.Errorf(ErrorSameEpisodeDownloadExists)
	}
	downloadingJobMap[job.key()] = job
	return job, nil
}

func (job *DownloadJob) startToDownload() error {
	log.Printf("Start to download: %s", job.key())
	torClient, err := torrent.NewClient(job.clientCfg)
	if err != nil {
		return err
	}
	tor, _ := torClient.AddMagnet(job.item.Link)
	<-tor.GotInfo()
	tor.DownloadAll()
	torClient.WaitAll()
	torClient.Close()
	log.Printf("completed downloading file: %s", job.key())

	for _, f := range tor.Files() {
		srcPath := path.Join(DownloadTmpFolder, f.Path())
		err = MoveFile(srcPath, job.destFolder, f.Path())
		if err != nil {
			return fmt.Errorf("failed to move downloaded file, cause: %s", err.Error())
		}
	}
	return nil
}

func MoveFile(srcPath, downloadFolder, episodeFolder string) error {
	inputFile, err := os.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("src file: %s does not exist", srcPath)
		}
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(srcPath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}