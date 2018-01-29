package logic

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type boltDao struct {
	Dao

	db *bolt.DB
}

var (
	// buckets
	bucketMeta  = []byte("metadata")
	bucketUsers = []byte("users")

	// constants
	versionName  = []byte("version")
	versionValue = []byte("perfcomp-1.0")
)

// NewBoltDao creates Bolt DB-based DAO
func NewBoltDao(dbPath string) (Dao, error) {
	var err error
	var result boltDao

	if result.db, err = bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 2 * time.Second, NoGrowSync: false}); err != nil {
		return nil, fmt.Errorf("unable to open DB: %v", err)
	}

	if err = result.db.Update(func(tx *bolt.Tx) error {
		meta := tx.Bucket(bucketMeta)
		if meta == nil {
			log.Printf("perform schema initialization for db=%s", dbPath)
			if meta, err = tx.CreateBucket(bucketMeta); err != nil {
				return fmt.Errorf("unable to create version bucket: %v", err)
			}

			if err = meta.Put(versionName, versionValue); err != nil {
				return err
			}

			if _, err = tx.CreateBucket(bucketUsers); err != nil {
				return fmt.Errorf("unable to create users bucket: %v", err)
			}
		} else {
			log.Printf("perform schema validation for db=%s", dbPath)
			actualVersionValue := meta.Get(versionName)

			if actualVersionValue == nil || !bytes.Equal(versionValue, actualVersionValue) {
				return fmt.Errorf("version mismatch, expected: %s, actual: %s", versionValue, actualVersionValue)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to perform initialization: %v", err)
	}

	return &result, nil
}

func (t *boltDao) Close() error {
	return t.db.Close()
}

func (t *boltDao) Add(profiles []*UserProfile) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(bucketUsers)
		if users == nil {
			return fmt.Errorf("unable to add profiles slice: users bucket is missing; data corrupted?")
		}

		for _, p := range profiles {
			var valueBuf bytes.Buffer
			encoder := gob.NewEncoder(&valueBuf)
			if err := encoder.Encode(p); err != nil {
				return fmt.Errorf("unable to encode profile=%s, error: %v", p, err)
			}

			if err := users.Put(getBytesFromID(p.ID), valueBuf.Bytes()); err != nil {
				return fmt.Errorf("unable to add profile=%s, error: %v", p, err)
			}
		}

		return nil
	})
}

func (t *boltDao) QueryUsers(offsetToken string, limit int) (*UserPage, error) {
	var result UserPage

	if err := t.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(bucketUsers)
		if users == nil {
			return fmt.Errorf("unable to query users: users bucket is missing; data corrupted?")
		}

		var k, v []byte

		cur := users.Cursor()
		if len(offsetToken) > 0 {
			keyBytes, err := hex.DecodeString(offsetToken)
			if err != nil {
				return fmt.Errorf("corrupted offset token: %v", err)
			}
			k, v = cur.Seek(keyBytes)
		} else {
			k, v = cur.First()
		}

		size := 0
		for {
			if k == nil {
				break
			}

			decoder := gob.NewDecoder(bytes.NewBuffer(v))
			var p UserProfile
			if err := decoder.Decode(&p); err != nil {
				return fmt.Errorf("unable to decode user profile value: offset=%d, offsetToken=%s, error=%v", size, offsetToken, err)
			}

			result.Profiles = append(result.Profiles, &p)

			k, v = cur.Next()

			size++
			if size > limit {
				if k != nil {
					result.OffsetToken = hex.EncodeToString(k)
				}
				break
			}
		}

		return nil

	}); err != nil {
		return nil, fmt.Errorf("unable to query users: %v", err)
	}

	return &result, nil
}

//
// Private
//

func getBytesFromID(id int) []byte {
	key := make([]byte, 4)
	binary.BigEndian.PutUint32(key, uint32(id))
	return key
}
