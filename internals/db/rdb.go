package db

import (
	// "bytes"
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

func InitRDBTrackers(conf *config.Config) {
	for _, rdb := range conf.RDB {
		tracker := NewSnapshotTracker(&rdb)
		trackers = append(trackers, tracker)

		go func() {
			defer tracker.ticker.Stop()

			for range tracker.ticker.C {
				// log.Printf("keys changed: %d - keys required to change: %d", tracker.keys, tracker.rdb.KeysChanged)
				if tracker.keys >= tracker.rdb.KeysChanged {
					SaveRDB(conf)
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

func SaveRDB(conf *config.Config) {
	fp := path.Join(conf.Dir, conf.RDBfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644) // owner (read-write), everyone else (read)
	if err != nil {
		log.Println("error opening rdb file: ", err)
		return
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(&DB.Store)
	if err != nil {
		log.Println("error saving to RDB: ", err)
		return
	}
	log.Println("saving DB to RDB file")
	// var buf bytes.Buffer
	// if state.bgsaveRunning {
	// 	err = gob.NewEncoder(&buf).Encode(&state.dbCopy)
	// } else {
	// 	DB.Mu.RLock()
	// 	err = gob.NewEncoder(&buf).Encode(&DB.Store)
	// 	DB.Mu.RUnlock()
	// }

	// if err != nil {
	// 	log.Println("error encoding db: ", err)
	// 	return
	// }

	// data := buf.Bytes()

	// bsum, err := Hash(&buf)
	// if err != nil {
	// 	log.Println("rdb - cannot compute buf checksum: ", err)
	// 	return
	// }

	// _, err = file.Write(data)
	// if err != nil {
	// 	log.Println("rdb - cannot write to file: ", err)
	// 	return
	// }
	// if err := file.Sync(); err != nil {
	// 	log.Println("rdb - cannot flush file to disk: ", err)
	// 	return
	// }
	// if _, err := file.Seek(0, io.SeekStart); err != nil {
	// 	log.Println("rdb - cannot seek file: ", err)
	// 	return
	// }

	// fsum, err := Hash(file)
	// if err != nil {
	// 	log.Println("rdb - cannot compute file checksum: ", err)
	// 	return
	// }

	// if bsum != fsum {
	// 	log.Printf("rdb - buf and file checksums do not match:\nf=%s\nb=%s\n", fsum, bsum)
	// 	return
	// }

	// log.Println("saved RDB file")

	// state.rdbStats.rdb_last_save_ts = time.Now().Unix()
	// state.rdbStats.rdb_saves++
}

func SyncRDB(conf *config.Config) {
	fp := path.Join(conf.Dir, conf.RDBfn)
	file, err := os.OpenFile(fp, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Println("error opening rdb file: ", err)
		file.Close()
		return
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&DB.Store)
	if err != nil {
		log.Println("error decoding rdb file: ", err)
		return
	}
	log.Println("synced RDB")
}

func Hash(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}