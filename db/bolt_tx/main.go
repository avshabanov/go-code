package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func disposeTempDB(db *bolt.DB) {
	path := db.Path()

	if err := db.Close(); err != nil {
		log.Fatalf("Error while closing temp DB: %v", err)
	}

	if err := os.Remove(path); err != nil {
		log.Fatalf("Unable to remove file for temp DB: %v", err)
	}
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := ioutil.TempFile("", "boltdb-txtest-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}

var helloBucketName = []byte("Hello")

func tx(msg string, db *bolt.DB, q chan string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(helloBucketName)
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}

		time.Sleep(time.Second)

		b.Get([]byte("1"))

		if err := b.Put([]byte("1"), []byte("val-"+msg)); err != nil {
			return err
		}

		log.Printf("[%s] Inserted record", msg)

		time.Sleep(time.Second)

		return nil
	})
	if err != nil {
		log.Printf("[%s] Error: %v", msg, err)
	}

	q <- fmt.Sprintf("<done: %s>", msg)
}

func main() {
	db, err := bolt.Open(tempfile(), 0664, nil)
	if err != nil {
		panic(err)
	}
	defer disposeTempDB(db)

	fmt.Printf("Opened test DB: %s\n", db.Path())

	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(helloBucketName)
		if err != nil {
			return err
		}

		if err := b.Put([]byte("1"), []byte("initial")); err != nil {
			return err
		}

		log.Println("inserted test data")
		return nil
	}); err != nil {
		panic(err)
	}

	const count = 3
	execQueue := make(chan string, count)

	for i := 0; i < count; i++ {
		go tx(fmt.Sprintf("tx-%d", i+1), db, execQueue)
	}

	for i := 0; i < count; i++ {
		msg := <-execQueue
		fmt.Printf("msg %d = %s\n", i, msg)
	}
}
