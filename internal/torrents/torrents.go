package torrents

import (
	"context"
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

// ReplyFunc is a function that sends a reply to the user.
type ReplyFunc func(int64, string)

type DownloadRequest struct {
	MagnetLink string
	Category   string
	ChatId     int64
}

func (request DownloadRequest) ToString() string {
	return fmt.Sprintf("magnet [%s]", request.Category)
}

type DownloadListEntry struct {
	Name      string
	TorrentId string
	Category  string
}

type Downloader struct {
	config  *config.Config
	session *torrent.Session
	// Stores the mapping between torrent ID and the download request.
	downloads map[string]*DownloadRequest
	reply     ReplyFunc
	mutex     sync.Mutex
}

func NewDownloader(ctx context.Context, cfg *config.Config, reply ReplyFunc) (*Downloader, error) {
	config := torrent.DefaultConfig
	config.DataDir = cfg.WorkDir
	config.Database = cfg.DatabasePath
	config.FilePermissions = 0o755

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
		config:    cfg,
		session:   session,
		downloads: make(map[string]*DownloadRequest),
		reply:     reply,
	}
	go d.RunCleanupLoop(ctx)
	return &d, nil
}

func (d *Downloader) Shutdown() {
	d.session.Close()
}

func (d *Downloader) RunCleanupLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Termination signal received, shutting down the cleanup loop")
			return
		case <-time.After(time.Second * 5):
			err := d.Cleanup()
			if err != nil {
				log.Printf("Cleanup error: %s", err.Error())
			}
		}
	}
}

func (d *Downloader) Add(req *DownloadRequest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.reply(req.ChatId, "Starting download of "+req.ToString())

	torr, err := d.session.AddURI(req.MagnetLink, nil)
	if err != nil {
		return fmt.Errorf("could not add torrent to session: %w", err)
	}
	d.downloads[torr.ID()] = req
	return nil
}

func (d *Downloader) Cleanup() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	torrents := d.session.ListTorrents()
	for _, torr := range torrents {
		stats := torr.Stats()
		if stats.Status == torrent.Seeding || stats.Status == torrent.Stopped {
			log.Printf("Removing completed torrent %s", torr.Name())

			category := config.UnsortedCategory
			req, found := d.downloads[torr.ID()]
			if found {
				d.reply(req.ChatId, fmt.Sprintf("Download of [%s] %s completed", req.Category, torr.Name()))
				log.Printf("Found download request for torrent %s, category %s", torr.ID(), req.Category)
				category = req.Category
			} else {
				// TODO: Store the requests in a persistent storage to recover from restarts.
				log.Printf("Could not find download request for torrent %s", torr.Name())
			}

			targetDir := d.GetTargetDir(category)
			err := d.MoveDownloadedFiles(torr.RootDirectory(), targetDir)
			if err != nil {
				log.Printf("Could not move downloaded files: %s", err.Error())
				continue
			}

			log.Printf("Removing torrent %s from session", torr.ID())
			err = d.session.RemoveTorrent(torr.ID())
			if err != nil {
				log.Printf("could not remove torrent from session: %s", err)
				continue
			}
		}
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

// SafeMkdir must be called under d.mutex.
func (d *Downloader) SafeMkdir(parentDir string, desiredName string) (string, error) {
	newPath, err := d.NewPath(parentDir, desiredName)
	if err == nil {
		err = os.Mkdir(newPath, 0o755)
	}
	if err != nil {
		return "", err
	}
	return newPath, nil
}

// SafeMove must be called under d.mutex.
func (d *Downloader) SafeMove(src string, destDir string) (string, error) {
	log.Printf("Moving %s to %s", src, destDir)
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
			return fmt.Errorf("could not create directory %s: %w", destDir, err)
		}
	}
	for _, entry := range entriesToMove {
		src := path.Join(srcDir, entry)
		_, err := d.SafeMove(src, destDir)
		if err != nil {
			return fmt.Errorf("could not move %s to %s: %w", src, destDir, err)
		}
	}
	return nil
}

func (d *Downloader) StatusString() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()

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

	return status
}

func (d *Downloader) List() []DownloadListEntry {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	torrents := d.session.ListTorrents()
	entries := make([]DownloadListEntry, len(torrents))
	for i, torr := range d.session.ListTorrents() {
		category := config.UnsortedCategory
		req, found := d.downloads[torr.ID()]
		if found {
			category = req.Category
		}
		entries[i] = DownloadListEntry{
			Name:      torr.Name(),
			TorrentId: torr.ID(),
			Category:  category,
		}
	}
	return entries
}

func (d *Downloader) Cancel(torrentId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	err := d.session.RemoveTorrent(torrentId)
	if err != nil {
		return fmt.Errorf("could not remove torrent %s: %w", torrentId, err)
	}
	return nil
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
