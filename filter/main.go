package filter

import (
	"fmt"
	"os"
	// "bufio"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"encoding/csv"
	"io/ioutil"
	"strconv"
	"log"
	"sort"
	"time"
)
			
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
	Distance float64
}

func (b Body) hasInterestingMaterial() bool {
	desired := [6]string{"Ruthenium", "Antimony", "Yttrium", "Technetium", "Polonium", "Tellurium"}
	for _,m := range desired {
		_, ex := b.Materials[m]
		if ex {
			return true
		}
	}
	return false
}

type Candidates []Body

func (c Candidates) Len() int {
	return len(c)
}

func (c Candidates) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Candidates) Less(i, j int) bool {
	a := c[i]; b := c[j]
	if a.Distance == b.Distance {
		// if equidistant to reference point
		if a.DistanceToArrival == b.DistanceToArrival {
			// if equidistant to arrival star
			return a.Gravity > b.Gravity
		} else {
			return a.DistanceToArrival < b.DistanceToArrival
		}
	} else {
		return a.Distance < b.Distance
	}
}

func (c Candidates) WriteToJSON(fn string) {
	fileWriter, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWriter.Close()
	gzipWriter := gzip.NewWriter(fileWriter)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipWriter.Close()
	enc := json.NewEncoder(gzipWriter)
	err = enc.Encode(c)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Candidates) WriteToCSV(fn string) {
	fileWriter, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWriter.Close()
	enc := csv.NewWriter(fileWriter)
	headers := []string{"Name", "System", "Distance", "DistanceToArrival", "Gravity", "Temperature"}
	err = enc.Write(headers)
	if err != nil {
		log.Fatal(err)
	}
	for _,cand := range c {
		fields := []string{
			cand.Name,
			cand.SystemName,
			strconv.FormatFloat(cand.Distance, 'f', 2, 64),
			strconv.FormatFloat(cand.DistanceToArrival, 'f', 2, 64),
			strconv.FormatFloat(cand.Gravity, 'f', 2, 64),
			strconv.FormatFloat(cand.SurfaceTemperature, 'f', 1, 64),
			cand.VolcanismType,
		}
		err := enc.Write(fields)
		if err != nil {
			log.Fatal(err)
		}
	}
	enc.Flush()
}

func loadDistances(fn string) map[int64]float64 {
	fmt.Printf("loading distances from %s... ", fn)
	gobfile, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer gobfile.Close()

	dec := gob.NewDecoder(gobfile)
	distances := make(map[int64]float64)
	dec.Decode(&distances)
	fmt.Printf("done\n")

	return distances
}

func loadDump(fn string) Candidates {
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
	data, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		log.Fatal(err)
	}
	var bodies Candidates
	err = json.Unmarshal(data, &bodies)
	fmt.Printf("loaded %d bodies from dump\n", len(bodies))
	return bodies
}


// func findCandidates(fn string, distances map[int64]float64) map[int64]Body {
func findCandidates(fn string, distances map[int64]float64, limit int64) Candidates {
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
		log.Fatal("Error reading JSON array start:", err)
	}
	// fmt.Printf("%T: %v\n", t, t)
	var i int64
	i = 0
	bodies := make(Candidates, 0, 1200000)
	// while there's more, decode it
	for dec.More() {
		var body Body
		err := dec.Decode(&body)
		if err != nil {
			log.Fatal("Error reading JSON obj.:",err)
		}
		distance, found := distances[body.SystemId]
		if found &&
		body.IsLandable && 
		body.DistanceToArrival > 12000 && 
		body.VolcanismType != "No volcanism" && 
		body.hasInterestingMaterial() {
			body.Distance = distance
			bodies = append(bodies, body)
		}
		i++
		fmt.Printf("\rProcessed: %d / Found %d", i, len(bodies))
		if limit > 0 && i >= limit {
			break
		}
	}
	fmt.Println()
	// read array closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal("Error reading JSON array end:", err)
	}
	return bodies
}

func FilterBodies(bodies, gob, outfn, outfmt, reexp string, limit int64) {
	var candidates Candidates
	if reexp == "" {
		distances_start := time.Now()
		distances := loadDistances(gob)
		distances_time := time.Since(distances_start)
		fmt.Printf("loadDistances %s\n", distances_time)
	
		candidates_start := time.Now()
		candidates = findCandidates(bodies, distances, limit)
		candidates_time := time.Since(candidates_start)
		fmt.Printf("findCandidates %s\n", candidates_time)
	
		// bufio.NewReader(os.Stdin).ReadString('\n')
	
		sort_start := time.Now()
		sort.Sort(Candidates(candidates))
		sort_time := time.Since(sort_start)
		fmt.Printf("sort %s\n", sort_time)
	} else {
		load_start := time.Now()
		candidates = loadDump(reexp)
		load_time := time.Since(load_start)
		fmt.Printf("load %s\n", load_time)
	}
	
	switch outfmt {
	case "json":
		json_start := time.Now()
		candidates.WriteToJSON(outfn)
		json_time := time.Since(json_start)
		fmt.Printf("json %s\n", json_time)
	case "csv":
		csv_start := time.Now()
		candidates.WriteToCSV(outfn)
		csv_time := time.Since(csv_start)
		fmt.Printf("csv %s\n", csv_time)
	default:
		log.Fatalf("Invalid format: %s Must be one of [json,csv]", outfmt)
	}
}
