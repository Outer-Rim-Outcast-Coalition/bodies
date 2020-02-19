package main

import (
	"fmt"
	"os"
	"encoding/binary"
	"bytes"
	"log"
	bolt "go.etcd.io/bbolt"
)

type Coords struct {
	X, Y, Z float64
}

type System struct {
	Id, Id64 int64
	Name string
	Coords Coords
	Date string
}

// type Persister interface {
// 	Write(ds map[int64]float64)
// }

// type BoltPersister struct {
// 	BoltDB stuff
// }

// func (db *BoltPersister) Write(ds map[int]float64) {

// }


func main() {

	// get filename to read
	filename := os.Args[1]
	fmt.Printf("Opening file %s\n", filename)
	db, err := bolt.Open("bodies.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
		err = db.View(func (tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Distances"))
			var k int64
			k = 194580063412
			key := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(key, k)
			val := b.Get(key)
			valbuf := bytes.NewReader(val)
			var v float64
			binary.Read(valbuf, binary.LittleEndian, &v)
			fmt.Printf("%d: %.2f\n", k, v)
			return nil
		})
}
