package download

import (
	"bytes"
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/larryrun/tor-poller/pkg/config"
	"github.com/larryrun/tor-poller/pkg/mags"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
)

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

func NewDownload(conf *config.Config, item *mags.TargetItem, destFolder string) error {
	jobMapMutex.Lock()
	defer jobMapMutex.Unlock()

	if len(downloadingJobMap) >= int(conf.ConcurrentDownload) {
		return fmt.Errorf(ErrorMaxConcurrentJobReached)
	}
	_, ok := downloadingJobMap[item.Name+item.Episode.String()]
	if ok {
		return fmt.Errorf(ErrorSameEpisodeDownloadExists)
	}
	job, err := createNewJob(conf, item, destFolder)
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

func createNewJob(conf *config.Config, item *mags.TargetItem, destFolder string) (*DownloadJob, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = conf.TmpFolder
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
	for _, f := range tor.Files() {
		absPath, err := filepath.Abs(filepath.Join(job.clientCfg.DataDir, f.Path()))
		if err != nil {
			panic(err)
		}
		log.Printf("Downloading file: %s", absPath)
	}
	tor.DownloadAll()
	torClient.WaitAll()
	torClient.Close()
	log.Printf("completed downloading file: %s", job.key())

	for _, f := range tor.Files() {
		err = MoveTempFileToDest(job.clientCfg.DataDir, f.Path(), job.destFolder)
		if err != nil {
			return fmt.Errorf("failed to move downloaded file, cause: %s", err.Error())
		}
	}
	return nil
}

func MoveTempFileToDest(tempFolder, tempFilePath, destFolder string) error {
	tempFileAbsPath, err := filepath.Abs(filepath.Join(tempFolder, tempFilePath))
	if err != nil {
		return fmt.Errorf("failed to get the abs path of the tempFile: %s", tempFileAbsPath)
	}
	tempFileInfo, err := os.Stat(tempFileAbsPath)
	if err != nil {
		return fmt.Errorf("failed to read tempFileAbsPath, error: %v", err)
	}
	if tempFileInfo.IsDir() {
		tempFileFolderPath := filepath.Dir(tempFilePath)
		destFileFolderPath := filepath.Join(destFolder, tempFileFolderPath)
		//this means the downloaded item is a folder, we need to make sure the folder exists
		if err := os.MkdirAll(destFileFolderPath, 0777); err != nil {
			return fmt.Errorf("failed to create item folder: %s, err: %s", destFileFolderPath, err.Error())
		}
		log.Printf("Moving from %s to %s", tempFileAbsPath, destFileFolderPath)
		cpCmd := exec.Command("cp", "-R", tempFileAbsPath, destFileFolderPath)
		var out bytes.Buffer
		cpCmd.Stdout = &out
		if err := cpCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy tempFolder %s to dest folder: %s, error: %v, output: %s", tempFileAbsPath, destFileFolderPath, err, out.String())
		}
		log.Printf("File moved, deleting the temp file: %s", tempFileAbsPath)
		if err := os.RemoveAll(tempFileAbsPath); err != nil {
			log.Printf("failed to remove tempFile: %s, cause: %s", tempFileAbsPath, err.Error())
		}
		log.Printf("Temp file removed")

		tempFileFolderAbsPath, err := filepath.Abs(path.Join(tempFolder, tempFileFolderPath))
		if err != nil {
			log.Printf("failed to get the abs path of the tempFileFolder: %s", tempFileFolderPath)
		}
		tempFileFolderFileInfos, err := ioutil.ReadDir(tempFileFolderAbsPath)
		if err != nil {
			log.Printf("failed to read tempFileFolderInfo: %s, cause: %v", tempFileFolderAbsPath, err)
			return nil
		}
		if len(tempFileFolderFileInfos) == 0 {
			log.Printf("TempFileFolder %s is empty now, removing it", tempFileFolderAbsPath)
			err := os.Remove(tempFileFolderAbsPath)
			if err != nil {
				log.Printf("failed to remove tempFileFolder: %s, cause: %v", tempFileFolderAbsPath, err)
			}
		}
	} else {
		downloadedFileDestPath := path.Join(destFolder, tempFilePath)
		downloadedFileFolderPath := path.Dir(downloadedFileDestPath);
		var out bytes.Buffer
		mkdirCmd := exec.Command("mkdir", "-p", downloadedFileFolderPath)
		mkdirCmd.Stdout = &out
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("failed to mkdir for dest file %s, error: %v, output: %s", downloadedFileDestPath, err, out.String())
		}

		out = bytes.Buffer{}
		cpCmd := exec.Command("cp", tempFileAbsPath, downloadedFileDestPath)
		cpCmd.Stdout = &out
		if err := cpCmd.Run(); err != nil {
			return fmt.Errorf("failed to cp tempFile %s to dest path: %s, error: %v, output: %s", tempFileAbsPath, downloadedFileDestPath, err, out.String())
		}
		if err := os.RemoveAll(tempFileAbsPath); err != nil {
			log.Printf("failed to remove tempFile: %s, cause: %s", tempFileAbsPath, err.Error())
		}
	}
	return nil
}

