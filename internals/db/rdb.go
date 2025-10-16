package db

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/shivakuppa/Go_Redis/config"
)

type SnapshotTracker struct {
	keys   int
	ticker time.Ticker
	rdb    *config.RDBSnapshot
}

func NewSnapshotTracker(rdb *config.RDBSnapshot) *SnapshotTracker {
	return &SnapshotTracker{
		keys:   0,
		ticker: *time.NewTicker(time.Second * time.Duration(rdb.Secs)),
		rdb:    rdb,
	}
}

var trackers = []*SnapshotTracker{}

func InitRDBTrackers(state *AppState) {
	for _, rdb := range state.Config.RDB {
		tracker := NewSnapshotTracker(&rdb)
		trackers = append(trackers, tracker)

		go func() {
			defer tracker.ticker.Stop()

			for range tracker.ticker.C {
				// log.Printf("keys changed: %d - keys required to change: %d", tracker.keys, tracker.rdb.KeysChanged)
				if tracker.keys >= tracker.rdb.KeysChanged {
					SaveRDB(state)
				}
				tracker.keys = 0
			}
		}()
	}
}

func IncrRDBTrackers() {
	for _, t := range trackers {
		t.keys++
	}
}

func SaveRDB(state *AppState) {
	fp := path.Join(state.Config.Dir, state.Config.RDBfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("error opening rdb file:", err)
		return
	}
	defer file.Close()

	log.Println("saving DB to RDB file")

	var buf bytes.Buffer
	var encodeErr error

	if state.BgSaveRunning {
		encodeErr = gob.NewEncoder(&buf).Encode(&state.DBCopy)
	} else {
		DB.mu.RLock()
		encodeErr = gob.NewEncoder(&buf).Encode(DB.GetItems())
		DB.mu.RUnlock()
	}

	if encodeErr != nil {
		log.Println("error encoding db:", encodeErr)
		return
	}

	data := buf.Bytes()

	// Compute hash of buffer before writing
	bsum, err := Hash(bytes.NewReader(data))
	if err != nil {
		log.Println("rdb - cannot compute buf checksum:", err)
		return
	}

	// Write to file
	if _, err := file.Write(data); err != nil {
		log.Println("rdb - cannot write to file:", err)
		return
	}
	if err := file.Sync(); err != nil {
		log.Println("rdb - cannot flush file to disk:", err)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("rdb - cannot seek file:", err)
		return
	}

	fsum, err := Hash(file)
	if err != nil {
		log.Println("rdb - cannot compute file checksum:", err)
		return
	}

	if bsum != fsum {
		log.Printf("rdb - buf and file checksums do not match:\nf=%s\nb=%s\n", fsum, bsum)
		return
	}

	log.Println("saved RDB file successfully!")
}


func SyncRDB(state *AppState) {
	fp := path.Join(state.Config.Dir, state.Config.RDBfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Println("error opening rdb file: ", err)
		file.Close()
		return
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&DB.store)
	if err != nil {
		log.Println("error decoding rdb file: ", err)
		return
	}
	log.Println("synced RDB")
}

func Hash(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
