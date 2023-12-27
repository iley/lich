package torrents

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/rain/torrent"
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
	requests chan *DownloadRequest
	config   *config.Config
	session  *torrent.Session
	mutex    sync.Mutex
}

func NewDownloader(ctx context.Context, cfg *config.Config) (*Downloader, error) {
	config := torrent.DefaultConfig
	config.DataDir = cfg.WorkDir
	config.Database = cfg.DatabasePath
	session, err := torrent.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("could not create torrent session: %w", err)
	}

	log.Println("Cleaning up state from previous runs")
	torrents := session.ListTorrents()
	for _, torrent := range torrents {
		log.Printf("Removing leftover torrent %s", torrent.ID())
		err = session.RemoveTorrent(torrent.ID())
		if err != nil {
			return nil, fmt.Errorf("could not remove torrent %s: %w", torrent.ID(), err)
		}
	}

	log.Println("Cleaning up work directory")
	err = removeDirectoryContents(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("could not clean up work directory %s: %w", cfg.WorkDir, err)
	}

	d := Downloader{
		requests: make(chan *DownloadRequest, 32),
		config:   cfg,
		session:  session,
	}
	go d.RunDownloadLoop(ctx)
	return &d, nil
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

func (d *Downloader) RunDownloadLoop(ctx context.Context) {
	defer d.session.Close()

	for {
		select {
		case <-ctx.Done():
			log.Print("Shutting down download loop")
			return
		case request := <-d.requests:
			request.Reply("Starting download of " + request.ToString())
			replyText := ""
			torrentName, err := d.Download(request)
			if err == nil {
				replyText = "Finished download of " + torrentName
			} else {
				replyText = fmt.Sprintf("Failed download of %s: %s", request.ToString(), err.Error())
			}
			log.Println(replyText)
			request.Reply(replyText)
		}
	}
}

func (d *Downloader) Download(request *DownloadRequest) (string, error) {
	torrent, err := d.session.AddURI(request.MagnetLink, nil)
	if err != nil {
		return "", fmt.Errorf("could not add torrent to session: %w", err)
	}

	completeChan := torrent.NotifyComplete()
	<-completeChan

	stopChan := torrent.NotifyStop()
	torrent.Stop()
	<-stopChan

	targetDir := d.GetTargetDir(request.Category)
	err = d.MoveDownloadedFiles(torrent.RootDirectory(), targetDir)
	if err != nil {
		return "", fmt.Errorf("could not move downloaded files: %w", err)
	}

	torrentName := torrent.Name()

	err = d.session.RemoveTorrent(torrent.ID())
	if err != nil {
		return "", fmt.Errorf("could not remove torrent from session: %w", err)
	}
	return torrentName, nil
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

func (d *Downloader) StatusString() string {
	status := ""

	torrents := d.session.ListTorrents()
	torrentStatuses := make([]string, len(torrents))
	for i, torrent := range d.session.ListTorrents() {
		stats := torrent.Stats()
		torrentStatuses[i] = fmt.Sprintf("%d: %s (%s)", i+1, torrent.Name(), stats.Status.String())
	}

	if len(torrentStatuses) > 0 {
		status += fmt.Sprintf("Active downloads:\n%s", strings.Join(torrentStatuses, "\n"))
	} else {
		status += "No active downloads"
	}

	if len(d.requests) > 0 {
		status += fmt.Sprintf("\nDownloads queued: %d", len(d.requests))
	} else {
		status += "\nNo downloads in the queue"
	}

	return status
}

func dateString() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

func removeDirectoryContents(dir string) error {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range dirEntries {
		err = os.RemoveAll(path.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
