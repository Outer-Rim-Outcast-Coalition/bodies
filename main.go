package main

import (
	"fmt"
	"os"
	"compress/gzip"
	"encoding/json"
	// "encoding/binary"
	// "bytes"
	"log"
	// "math"
	// bolt "go.etcd.io/bbolt"
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
			
type Body struct {
	Id, Id64, BodyId int64
	Name, Type, SubType string
	Parents []map[string]int
	DistanceToArrival float64
	IsLandable bool
	Gravity, EarthMasses, Radius float64
	SurfaceTemperature, SurfacePressure float64
	VolcanismType, AtmosphereType string
	AtmosphereComposition, SolidComposition map[string]float64
	TerraformingState string
	OrbitalPeriod, SemiMajorAxis, OrbitalEccentricity float64
	OrbitalInclination, ArgOfPeriapsis, RotationalPeriod float64
	RotationalPeriodTidallyLocked bool
	AxialTilt float64
	Materials map[string]float64
	UpdateTime string
	SystemId, SystemId64 int64
	SystemName string
}
func filterBodies(fn string) map[string]Body {
	bodies := make(map[string]Body)
	// open file, chain gzip.Reader and json.Decoder
	fileReader, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()
	gzipReader, err := gzip.NewReader(fileReader)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()
	dec := json.NewDecoder(gzipReader)
	if err != nil {
		log.Fatal(err)
	}

	// read array opening bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%T: %v\n", t, t)
	i := 0
	// while there's more, decode it
	for dec.More() {
		var m Body
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		bodies[m.Name] = m
		// fmt.Printf("\rProcessed: %d / Found: %d", i+1, )
		i++
		if i >= 20 {
			break
		}
	}
	fmt.Println()
	// read array closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	return bodies
}

// func batchWriter(id int, items <-chan distanceStruct) {
// 	for i := range items {

// 	}
// }

func main() {

	// get filename to read
	filename := os.Args[1]
	fmt.Printf("Opening file %s\n", filename)
	// distances := computeDistances(filename)
	// fmt.Printf("%s distances computed\n", len(distances))
	// db, err := bolt.Open("bodies.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// db.Update(func(tx *bolt.Tx) error {
	// 	_, err := tx.CreateBucket([]byte("Distances"))
	// 	if err != nil {
	// 		return fmt.Errorf("create bucket: %s", err)
	// 	}
	// 	return nil
	// })
	// writtenCount := 0
	// for writtenCount < len(distances) {
	// 	thisBatch := 0
	// 	err := db.Batch(func(tx *bolt.Tx) error {
	// 		b := tx.Bucket([]byte("Distances"))
	// 		for k, v := range distances {
	// 			key := make([]byte, binary.MaxVarintLen64)
	// 			binary.PutVarint(key, k)
	// 			val := new(bytes.Buffer)
	// 			err := binary.Write(val, binary.LittleEndian, v)
	// 			if err != nil {
	// 				fmt.Println("binary.Write failed:", err)
	// 			}
	// 			err = b.Put(key, val.Bytes())
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 			thisBatch++
	// 			if thisBatch >= 100000 {
	// 				writtenCount += thisBatch
	// 				break
	// 			}
	// 		}
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		fmt.Println("Batch write error:", err)
	// 	}
	// }
	bodies := filterBodies(filename)
	for name, sysdata := range bodies {
		fmt.Printf("%s: %#v\n", name, sysdata)
	}
}
