package download

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewDownload(t *testing.T) {
	//link := "magnet:?xt=urn:btih:2b8bf8d40337321d5385c8e224aab53841daf71c&dn=Shameless.US.S09E01.Are.You.There.Shim.Its.Me.Ian.1080p.AMZN.WEBRip.DDP5.1.x264-NTb%5Brartv%5D&tr=http%3A%2F%2Ftracker.trackerfix.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2710&tr=udp%3A%2F%2F9.rarbg.to%3A2710"
	//fileName := "Shameless.US.S09E01.Are.You.There.Shim.Its.Me.Ian.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv"
	//destFolder := "~"
}

func TestMoveTempFileToDest(t *testing.T) {
	destFolder := "./testdata/download_dest"
	downloadTemp := "./testdata/download_tmp"

	clearTestFolder := func(folder string) {
		infos, _ := ioutil.ReadDir(folder)
		for _, info := range infos {
			if info.Name() != ".keep" {
				os.RemoveAll(path.Join(folder, info.Name()))
			}
		}
	}
	defer clearTestFolder(destFolder)
	defer clearTestFolder(downloadTemp)

	//test moving a folder
	//create folder and file
	downloadedFolderName := "downloadedParentFolder"
	downloadedFolderPath := path.Join(downloadTemp, downloadedFolderName)
	err := os.Mkdir(downloadedFolderPath, 0777)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	downloadedFileName := "downloaded.txt"
	downloadedFilePath := path.Join(downloadedFolderPath, downloadedFileName)
	_, err = os.Create(downloadedFilePath)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	testContent := "testContent"
	err = ioutil.WriteFile(downloadedFilePath, []byte(testContent), 0777)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	downloadedFileRelativePath := path.Join(downloadedFolderName, downloadedFileName)
	err = MoveTempFileToDest(downloadTemp, downloadedFileRelativePath, destFolder)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	actualContent, err := readFileContent(path.Join(destFolder, downloadedFileRelativePath))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, testContent, actualContent)

	//test moving a file inside a folder
	downloadedSubFolderName := "downloadedSubFolder"
	downloadedSubFolderPath := path.Join(downloadTemp, downloadedFolderName, downloadedSubFolderName)
	downloadedSubFilePath := path.Join(downloadedSubFolderPath, downloadedFileName)
	err = os.MkdirAll(downloadedSubFolderPath, 0777)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	_, err = os.Create(downloadedSubFilePath)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	err = ioutil.WriteFile(downloadedSubFilePath, []byte(testContent), 0777)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	downloadedSubFolderRelativePath := path.Join(downloadedFolderName, downloadedSubFolderName)
	err = MoveTempFileToDest(downloadTemp, downloadedSubFolderRelativePath, destFolder)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	actualContent, err = readFileContent(path.Join(destFolder, downloadedSubFolderRelativePath, downloadedFileName))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, testContent, actualContent)

	//test moving a file
	downloadedFilePath = path.Join(downloadTemp, downloadedFileName)
	_, err = os.Create(downloadedFilePath)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	err = ioutil.WriteFile(downloadedFilePath, []byte(testContent), 0777)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	downloadedFileRelativePath = downloadedFileName
	err = MoveTempFileToDest(downloadTemp, downloadedFileRelativePath, destFolder)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	actualContent, err = readFileContent(path.Join(destFolder, downloadedFileRelativePath))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, testContent, actualContent)
}

func readFileContent(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}