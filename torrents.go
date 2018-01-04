package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

type TorrentDownloader struct {
	requests chan *DownloadRequest
	config   *TorrentConfig
	mutex    sync.Mutex
}

type ReplyFunc func(string)

type DownloadRequest struct {
	MagnetLink string
	Category   string
	Reply      ReplyFunc
}

func NewTorrentDownloader(config *TorrentConfig) (*TorrentDownloader, error) {
	if err := ValidateTorrentConfig(config); err != nil {
		return nil, err
	}
	d := TorrentDownloader{
		requests: make(chan *DownloadRequest, 32),
		config:   config,
	}
	go d.RunDownloadLoop()
	return &d, nil
}

func ValidateTorrentConfig(config *TorrentConfig) error {
	if config.WorkDir == "" {
		return errors.New("Missing required option 'work_directory'")
	}
	err := os.MkdirAll(config.WorkDir, 0755)
	if err != nil {
		msg := fmt.Sprintf("Could not create work directory %s: %s",
			config.WorkDir, err.Error())
		return errors.New(msg)
	}
	hasUnsortedCategory := false
	for category, targetDir := range config.TargetDirs {
		if category == UnsortedCategory {
			hasUnsortedCategory = true
		}
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			msg := fmt.Sprintf("Could not create target directory %s for category %s: %s",
				targetDir, category, err.Error())
			return errors.New(msg)
		}
	}
	if !hasUnsortedCategory {
		return errors.New(fmt.Sprintf("Required category '%s' not found", UnsortedCategory))
	}
	return nil
}

// Must be called under d.mutex. Returns full path.
func (d *TorrentDownloader) NewPath(parentDir string, desiredName string) (string, error) {
	for index := 0; ; index++ {
		path := path.Join(parentDir, desiredName)
		if index > 0 {
			path = fmt.Sprintf("%s_%d", path, index)
		}
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return path, nil
		} else if err != nil {
			return "", err
		}
	}
}

func (d *TorrentDownloader) SafeMkdir(parentDir string, desiredName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	newPath, err := d.NewPath(parentDir, desiredName)
	if err == nil {
		err = os.Mkdir(newPath, 0755)
	}
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func (d *TorrentDownloader) SafeMove(src string, destDir string) (string, error) {
	log.Printf("Moving %s to %s", src, destDir)
	d.mutex.Lock()
	defer d.mutex.Unlock()
	base := path.Base(src)
	dest, err := d.NewPath(destDir, base)
	if err != nil {
		return "", err
	}
	err = os.Rename(src, dest)
	if err != nil {
		return "", err
	}
	return dest, nil
}

func (d *TorrentDownloader) RunDownloadLoop() {
	for request := range d.requests {
		request.Reply("Starting download of " + request.ToString())
		replyText := ""
		err := d.Download(request)
		if err == nil {
			replyText = "Finished download of " + request.ToString()
		} else {
			replyText = fmt.Sprintf("Failed download of %s: %s", request.ToString(), err.Error())
		}
		log.Println(replyText)
		request.Reply(replyText)
	}
}

func (d *TorrentDownloader) Download(request *DownloadRequest) error {
	workDir, err := d.SafeMkdir(d.config.WorkDir, DateString())
	if err != nil {
		msg := fmt.Sprintf("Could not create work directory %s: %s", workDir, err.Error())
		return errors.New(msg)
	}
	logFile := path.Join(workDir, "aria2.log")
	cmd := exec.Command("aria2c",
		"-d", workDir,
		"--log", logFile,
		"--log-level", "notice",
		"--seed-time", "0",
		request.MagnetLink)
	err = cmd.Run()
	if err != nil {
		return err
	}
	targetDir := d.GetTargetDir(request.Category)
	return d.MoveDownloadedFiles(workDir, targetDir)
}

func (d *TorrentDownloader) GetTargetDir(category string) string {
	targetDir, found := d.config.TargetDirs[category]
	if !found {
		return d.config.TargetDirs[UnsortedCategory]
	}
	return targetDir
}

func (d *TorrentDownloader) MoveDownloadedFiles(srcDir string, destDir string) error {
	fileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil
	}
	entriesToMove := make([]string, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		entry := fileInfo.Name()
		if strings.HasSuffix(entry, ".log") {
			continue
		}
		entriesToMove = append(entriesToMove, entry)
	}
	if len(entriesToMove) > 1 {
		destDir, err = d.SafeMkdir(destDir, "torrent")
		if err != nil {
			return err
		}
	}
	for _, entry := range entriesToMove {
		src := path.Join(srcDir, entry)
		_, err := d.SafeMove(src, destDir)
		if err != nil {
			msg := fmt.Sprintf("Could not move %s to %s: %s", src, destDir, err.Error())
			return errors.New(msg)
		}
	}
	return nil
}

func (d *TorrentDownloader) Categories() []string {
	categories := make([]string, 0, len(d.config.TargetDirs))
	for category, _ := range d.config.TargetDirs {
		categories = append(categories, category)
	}
	return categories
}

func DateString() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

func (request DownloadRequest) ToString() string {
	return fmt.Sprintf("magnet [%s]", request.Category)
}
