package filter

import (
	"fmt"
	"os"
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"sync"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GasGiantKey struct {
	SystemId, BodyId int64
}

type GasGiantValue struct{}

type Body struct {
	Id, Id64, BodyId                                     int64
	Name, Type, SubType                                  string
	Parents                                              []map[string]int
	DistanceToArrival                                    float64
	IsLandable                                           bool
	Gravity, EarthMasses, Radius                         float64
	SurfaceTemperature, SurfacePressure                  float64
	VolcanismType, AtmosphereType                        string
	AtmosphereComposition, SolidComposition              map[string]float64
	TerraformingState                                    string
	OrbitalPeriod, SemiMajorAxis, OrbitalEccentricity    float64
	OrbitalInclination, ArgOfPeriapsis, RotationalPeriod float64
	RotationalPeriodTidallyLocked                        bool
	AxialTilt                                            float64
	Materials                                            map[string]float64
	UpdateTime                                           string
	SystemId, SystemId64                                 int64
	SystemName                                           string
	Distance                                             float64
}

func (b Body) hasInterestingMaterial() bool {
	desired := [6]string{"Ruthenium", "Antimony", "Yttrium", "Technetium", "Polonium", "Tellurium"}
	for _, m := range desired {
		_, ex := b.Materials[m]
		if ex {
			return true
		}
	}
	return false
}

type Candidates []Body

type DistanceMap map[int64]float64

func (c Candidates) Len() int {
	return len(c)
}

func (c Candidates) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Candidates) Less(i, j int) bool {
	a := c[i]
	b := c[j]
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
	for _, cand := range c {
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

func loadDistances(fn string) DistanceMap {
	fmt.Printf("loading distances from %s... ", fn)
	gobfile, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer gobfile.Close()

	dec := gob.NewDecoder(gobfile)
	distances := make(DistanceMap)
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

func inputWorker(id int, queue <-chan string, results chan<- Body, ggs chan<- GasGiantKey, distances DistanceMap, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting...\n", id)
	ggRegEx, _ := regexp.Compile("(?i)gas giant")
	var k GasGiantKey
	for t := range queue {
		t = t[:len(t) - 1]
		body := Body{}
		json.Unmarshal([]byte(t), &body)
		distance, found := distances[body.SystemId]
		if found &&
		body.IsLandable &&
		body.DistanceToArrival > 12000 &&
		body.VolcanismType != "No volcanism" &&
		body.SurfaceTemperature <= 220 &&
		body.hasInterestingMaterial() {
			body.Distance = distance
			results <- body
		}
		if found && body.Type == "Planet" {
			if ggRegEx.MatchString(body.SubType) {
				// fmt.Println()
				k = GasGiantKey{SystemId: body.SystemId, BodyId: body.BodyId}
				ggs <- k
			}
		}
	}
	fmt.Printf("Worker %d done.\n", id)
}

func gasGiantsWorker(incoming <-chan GasGiantKey, gasGiants map[GasGiantKey]GasGiantValue, wg *sync.WaitGroup) {
	wg.Done()
	for in := range incoming {
		gasGiants[in] = struct{}{}
	}
}
func outputWorker(results <-chan Body, candidates *Candidates, wg *sync.WaitGroup) {
	defer wg.Done()
	for c := range results {
		*candidates = append(*candidates, c)
	}
}

// func findCandidates(fn string, distances map[int64]float64) map[int64]Body {
func findCandidates(fn string, distances DistanceMap, limit int64) Candidates {
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
	scanner := bufio.NewScanner(gzipReader)
	var i int64
	i = 0
	var k GasGiantKey
	gasGiants := make(map[GasGiantKey]GasGiantValue)
	var t string
	in_queue := make(chan string, 100000)
	results := make(chan Body, 100000)
	gasGiantsChan := make(chan GasGiantKey, 100000)
	candidates := make(Candidates, 0, 200000)
	var gg_wg sync.WaitGroup
	go gasGiantsWorker(gasGiantsChan, gasGiants, &gg_wg)
	gg_wg.Add(1)
	var in_wg sync.WaitGroup
	for w := 0; w < 3; w++ {
		go inputWorker(w, in_queue, results, gasGiantsChan, distances, &in_wg)
		in_wg.Add(1)
	}
	var out_wg sync.WaitGroup
	go outputWorker(results, &candidates, &out_wg)
	out_wg.Add(1)
	for scanner.Scan() {
		t = scanner.Text()
		if strings.HasSuffix(t, ",") {
			in_queue <- t
			i++
			fmt.Printf("\rRead: %d / Queued to process: %d / Queued to keep: %d / Candidates: %d", i, len(in_queue), len(results), len(candidates))
			if limit > 0 && i >= limit {
				break
			}
		}
	}
	fmt.Println()
	close(in_queue)
	in_wg.Wait()
	close(results)
	close(gasGiantsChan)
	out_wg.Wait()
	candidatesFinal := make(Candidates, 0, len(candidates))
	var c, f int
	c, f = 0, 0
	for _,body := range candidates {
		var pId int
		var isPlanet bool
		for _,p := range body.Parents {
			pId, isPlanet = p["Planet"]
			if isPlanet {
				k = GasGiantKey{SystemId: body.SystemId, BodyId: int64(pId)}
				break
			}
		}
		if isPlanet {
			_, isGasGiant := gasGiants[k]
			if isGasGiant {
				candidatesFinal = append(candidatesFinal, body)
				f++
			}
		}
		c++
		fmt.Printf("\rCandidates processed: %d / Orbiting gas giants: %d", c, f)
	}
	fmt.Println()

	return candidatesFinal
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

	bufio.NewReader(os.Stdin).ReadString('\n')

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
