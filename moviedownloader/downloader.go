package moviedownloader

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

const DownloadTmpFolder = "./download_tmp"
const MaxConcurrentDownloading = 5
const ErrorMaxConcurrentJobReached = "max concurrent job count reached"

var jobMapMutex sync.Mutex
var downloadingJobMap map[string]*DownloadJob

func init() {
	downloadingJobMap = make(map[string]*DownloadJob, 0)
}

type DownloadJob struct {
	clientCfg  *torrent.ClientConfig
	link       string
	fileName   string
	destFolder string
}

func NewDownload(link, fileName, destFolder string) error {
	jobMapMutex.Lock()
	if len(downloadingJobMap) >= MaxConcurrentDownloading {
		return fmt.Errorf(ErrorMaxConcurrentJobReached)
	}
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = DownloadTmpFolder
	job := &DownloadJob{link: link, fileName: fileName, destFolder: destFolder, clientCfg: cfg}
	downloadingJobMap[fileName] = job
	go func() {
		err := job.startToDownload()
		delete(downloadingJobMap, fileName)
		if err != nil {
			log.Printf("Failed to download: %s, cause: %s", fileName, err.Error())
		} else {
			log.Printf("%s has been downloaded successfully", fileName)
		}
	}()
	jobMapMutex.Unlock()
	return nil
}

func (job *DownloadJob) startToDownload() error {
	torClient, err := torrent.NewClient(job.clientCfg)
	defer torClient.Close()
	if err != nil {
		return err
	}
	tor, _ := torClient.AddMagnet(job.link)
	<-tor.GotInfo()
	tor.DownloadAll()
	completed := torClient.WaitAll()
	if completed {
		for _, f := range tor.Files() {
			srcPath := path.Join(DownloadTmpFolder, f.Path())
			folderName := job.fileName[:strings.LastIndex(job.fileName, ".")]
			destPath := path.Join(job.destFolder, folderName)
			err = MoveFile(srcPath, destPath)
			if err != nil {
				return fmt.Errorf("failed to move downloaded file, cause: %s", err.Error())
			}
		}
		return nil
	} else {
		return fmt.Errorf("download not complete")
	}
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
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
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}