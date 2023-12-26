package torrents

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/iley/lich/internal/config"
)

type ReplyFunc func(string)

type DownloadRequest struct {
	MagnetLink string
	Category   string
	Reply      ReplyFunc
}

func (request DownloadRequest) ToString() string {
	return fmt.Sprintf("magnet [%s]", request.Category)
}

type Downloader struct {
	requests        chan *DownloadRequest
	config          *config.Config
	inProgressCount int
	mutex           sync.Mutex
}

func NewDownloader(cfg *config.Config) (*Downloader, error) {
	d := Downloader{
		requests: make(chan *DownloadRequest, 32),
		config:   cfg,
	}
	go d.RunDownloadLoop()
	return &d, nil
}

func (d *Downloader) GetInProgressCount() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.inProgressCount
}

func (d *Downloader) IncrInProgressCount(incr int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.inProgressCount += incr
}

func (d *Downloader) AddRequest(req *DownloadRequest) error {
	select {
	case d.requests <- req:
	default:
		return errors.New("Download queue full")
	}
	return nil
}

// NewPath must be called under d.mutex. Returns full path.
func (d *Downloader) NewPath(parentDir string, desiredName string) (string, error) {
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

func (d *Downloader) SafeMkdir(parentDir string, desiredName string) (string, error) {
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

func (d *Downloader) SafeMove(src string, destDir string) (string, error) {
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

func (d *Downloader) RunDownloadLoop() {
	for request := range d.requests {
		request.Reply("Starting download of " + request.ToString())
		replyText := ""
		d.IncrInProgressCount(1)
		err := d.Download(request)
		d.IncrInProgressCount(-1)
		if err == nil {
			replyText = "Finished download of " + request.ToString()
		} else {
			replyText = fmt.Sprintf("Failed download of %s: %s", request.ToString(), err.Error())
		}
		log.Println(replyText)
		request.Reply(replyText)
	}
}

func (d *Downloader) Download(request *DownloadRequest) error {
	workDir, err := d.SafeMkdir(d.config.WorkDir, dateString())
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

func (d *Downloader) GetTargetDir(category string) string {
	targetDir, found := d.config.TargetDirs[category]
	if !found {
		return d.config.TargetDirs[config.UnsortedCategory]
	}
	return targetDir
}

func (d *Downloader) MoveDownloadedFiles(srcDir string, destDir string) error {
	fileInfos, err := os.ReadDir(srcDir)
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

func dateString() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}
