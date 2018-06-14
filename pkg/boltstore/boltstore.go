//
// Package boltstore an boltdb data store for paho.mqtt.golang.
//
package boltstore

import (
	"bytes"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/eclipse/paho.mqtt.golang/packets"
	log "github.com/sirupsen/logrus"
)

var defaultConfig *ConfigType

// BoltStore implements the paho.mqtt.golang 'store' interface to provide a true
// persistence, even across client failure.
type BoltStore struct {
	// read-write lock
	sync.RWMutex
	// path boltdb location
	path string
	// opened is opened flag
	opened bool
	// conf config instance
	conf *ConfigType
	// db boltdb instance
	db *bolt.DB
}

// NewBoltStoreWithConfig returns a pointer to a new instance of BoltStore, the instance
// is not initialized and ready to use until Open() has been called on it.
// You should provide boltdb path with config instance to this function.
func NewBoltStoreWithConfig(path string, conf *ConfigType) *BoltStore {
	store := &BoltStore{
		path:   path,
		opened: false,
		conf:   conf,
		db:     nil,
	}
	return store
}

// NewBoltStore returns a pointer to a new instance of BoltStore, the instance
// is not initialized and ready to use until Open() has been called on it.
func NewBoltStore(path string) *BoltStore {
	return NewBoltStoreWithConfig(path, defaultConfig)
}

// Open will allow the BoltStore to be used.
func (store *BoltStore) Open() {
	log.Debugln("enter Open()")
	store.Lock()
	defer store.Unlock()

	if store.path == "" {
		filepath, _ := os.Getwd()
		// use unix timestamp as filename
		filename := strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".db"
		store.path = path.Join(filepath, filename)
	}

	db, err := bolt.Open(store.path, 0600, &bolt.Options{Timeout: time.Duration(store.conf.BoltStore.OpenTimeout) * time.Second})
	if err != nil {
		// fatal with return
		log.WithField("error", err).Fatal(ErrOpenDB)
	}

	store.db = db
	store.opened = true
	log.Infoln("store is opened at", store.path)
}

// Put will put a message into the store, associated with the provided key value.
func (store *BoltStore) Put(key string, message packets.ControlPacket) {
	log.Debugln("enter Put()")
	store.Lock()
	defer store.Unlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return
	}

	// leverage ControlPacket's 'Write' function, now we got an io.writer 'buf'.
	var buf bytes.Buffer
	if err := message.Write(&buf); err != nil {
		log.WithField("error", err).Error(ErrWriteControlPacket)
		return
	}

	// put message to bolt bucket, we omit the error return intentionally
	store.db.Update(func(tx *bolt.Tx) error {
		// create or get the bolt bucket 'bucketname'
		bk, err := tx.CreateBucketIfNotExists([]byte(store.conf.BoltStore.BucketName))
		if err != nil {
			log.WithField("error", err).Error(ErrCreateBucket)
			return nil
		}

		// we can get ControlPacket's byte array from buf.Bytes()
		if err := bk.Put([]byte(key), buf.Bytes()); err != nil {
			log.WithField("error", err).Error(ErrPutValue)
			return nil
		}
		return nil
	})
}

// Get takes a key and looks in the store for a matching message
// returning either the message pointer or nil.
func (store *BoltStore) Get(key string) packets.ControlPacket {
	log.Debugln("enter Get()")
	store.RLock()
	defer store.RUnlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return nil
	}

	var packet packets.ControlPacket

	// get message from bolt bucket, we omit the error return intentionally
	store.db.View(func(tx *bolt.Tx) error {
		// get the bolt bucket 'bucketname'
		bk := tx.Bucket([]byte(store.conf.BoltStore.BucketName))
		if bk == nil {
			log.Warningln(ErrGetBucket)
			return nil
		}

		packetBytes := bk.Get([]byte(key))
		if packetBytes == nil {
			log.WithField("key", key).Error(ErrGetValue)
			return nil
		}

		// convert ControlPacket byte array to io.reader
		var err error
		packet, err = packets.ReadPacket(bytes.NewBuffer(packetBytes))
		if err != nil {
			log.WithField("error", err).Error(ErrReadControlPacket)
			return nil
		}
		return nil
	})

	return packet
}

// All returns a slice of strings containing all the keys currently in the BoltStore.
func (store *BoltStore) All() []string {
	log.Debugln("enter All()")
	store.RLock()
	defer store.RUnlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return nil
	}

	keys := []string{}

	// get all keys from the bolt bucket, we omit the error return intentionally
	store.db.View(func(tx *bolt.Tx) error {
		// get the bolt bucket 'bucketname'
		bk := tx.Bucket([]byte(store.conf.BoltStore.BucketName))
		if bk == nil {
			log.Warningln(ErrGetBucket)
			return nil
		}

		// iterate all key value pairs
		bk.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
		return nil
	})

	return keys
}

// Del takes a key, searches the BoltStore and if the key is found
// deletes the message pointer associated with it.
func (store *BoltStore) Del(key string) {
	log.Debugln("enter Del()")
	store.Lock()
	defer store.Unlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return
	}

	// delete message from the bolt bucket, we omit the error return intentionally
	store.db.Update(func(tx *bolt.Tx) error {

		// get the bolt bucket 'bucketname'
		bk := tx.Bucket([]byte(store.conf.BoltStore.BucketName))
		if bk == nil {
			log.Warningln(ErrGetBucket)
			return nil
		}

		// delete message by key
		if err := bk.Delete([]byte(key)); err != nil {
			log.WithField("error", err).Error(ErrDelValue)
			return nil
		}
		return nil
	})
}

// Close will disallow modifications to the state of the store.
func (store *BoltStore) Close() {
	log.Debugln("enter Close()")
	store.Lock()
	defer store.Unlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return
	}

	// we don't check whether the store.db exists or not
	store.db.Close()
	store.opened = false
	log.Infoln("Bolt store closed")
}

// Reset eliminates all persisted message data in the store.
func (store *BoltStore) Reset() {
	log.Debugln("enter Reset()")
	store.Lock()
	defer store.Unlock()

	if !store.opened {
		log.Warningln(ErrUseDB)
		return
	}

	// delete the default bucket
	store.db.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket([]byte(store.conf.BoltStore.BucketName)); err != nil {
			log.WithField("error", err).Warn(ErrDelBucket)
			return nil
		}
		return nil
	})
	log.Infoln("Bolt store wiped")
}
