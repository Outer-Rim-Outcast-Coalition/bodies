package distances

import (
	"fmt"
	"os"
	"compress/gzip"
	"encoding/json"
	"encoding/gob"
	"log"
	"math"
)

type Coordinates struct {
	X, Y, Z float64
}

type System struct {
	Id, Id64 int64
	Name string
	Coords Coordinates
	Date string
}

type Distances map[int64]float64

func getDistances(fn string, max, min float64) Distances {
	distances := make(Distances)
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
	fmt.Printf("opened %s...\n", fn)
	// read array opening bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%T: %v\n", t, t)
	i := 0
	j := 0
	// while there's more, decode it
	for dec.More() {
		var s System
		err := dec.Decode(&s)
		if err != nil {
			log.Fatal(err)
		}
		d := math.Sqrt(math.Pow(s.Coords.X, 2) + math.Pow(s.Coords.Y, 2) + math.Pow(s.Coords.Z, 2))
		if d <= max && d > min {
			distances[s.Id] = d
			j++
		}
		fmt.Printf("\rProcessed: %d / Found: %-10d", i, j)
		i++
	}
	fmt.Println()
	// read array closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	return distances
}

func (d Distances) WriteToGob(fn string) {
	// create output file
	gobfile, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer gobfile.Close()
	// create encoder
	enc := gob.NewEncoder(gobfile)
	if err := enc.Encode(d); err != nil {
		log.Fatal(err)
	}
}

func MakeDB(infile, gobfile string, max, min float64) {
	distances := getDistances(infile, max, min)
	distances.WriteToGob(gobfile)
}