package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var helloBucketName = []byte("Hello")
var keyOne = []byte{1}

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

func updateCounter(msg string, key []byte, db *bolt.DB, q chan string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(helloBucketName)
		if b == nil {
			return fmt.Errorf("bucket does not exist")
		}

		time.Sleep(time.Second)

		val := b.Get(key)

		if err := b.Put(key, []byte(string(val)+":"+msg)); err != nil {
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

func demoOverlappingUpdates(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(helloBucketName)
		if err != nil {
			return err
		}

		if err := b.Put(keyOne, []byte("init-0")); err != nil {
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
		go updateCounter(fmt.Sprintf("tx-%d", i+1), keyOne, db, execQueue)
	}

	for i := 0; i < count; i++ {
		msg := <-execQueue
		fmt.Printf("msg %d = %s\n", i, msg)
	}

	// show result
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(helloBucketName)
		if b == nil {
			return fmt.Errorf("Hello bucket doesn't exist")
		}

		val := b.Get(keyOne)
		log.Printf("End result: 1 => %s", string(val))

		return nil
	})
}

func demoIndependentUpdates(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(helloBucketName)
		if err != nil {
			return err
		}

		if err := b.Put(keyOne, []byte("init-0")); err != nil {
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
		go updateCounter(fmt.Sprintf("tx-%d", i+1), []byte{byte(i + 1)}, db, execQueue)
	}

	for i := 0; i < count; i++ {
		msg := <-execQueue
		fmt.Printf("msg %d = %s\n", i, msg)
	}

	// show result
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(helloBucketName)
		if b == nil {
			return fmt.Errorf("Hello bucket doesn't exist")
		}

		val := b.Get(keyOne)
		log.Printf("End result: 1 => %s", string(val))

		return nil
	})
}

func main() {
	mode := flag.String("mode", "overlapping-updates", "Specifies what demo to run")
	flag.Parse()

	db, err := bolt.Open(tempfile(), 0664, nil)
	if err != nil {
		panic(err)
	}
	defer disposeTempDB(db)

	fmt.Printf("Opened test DB: %s\n", db.Path())

	switch *mode {
	case "overlapping-updates":
		demoOverlappingUpdates(db)
	case "independent-updates":
		demoIndependentUpdates(db)
	default:
		fmt.Printf("Wrong mode: %s\n", *mode)
		flag.Usage()
	}
}
