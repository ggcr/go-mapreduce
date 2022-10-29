package worker

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Output struct {
	mu      sync.Mutex
	out_map map[string]int
}

func decodeDict(filepath string) map[string]int {
	v, err := os.Open(filepath)
	check(err)

	var decodedMap map[string]int
	d := gob.NewDecoder(v)

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	// Ta da! It is a map!
	fmt.Printf("length for readed dict (Key=\"%s\") => %d\n", filepath, len(decodedMap))
	return decodedMap
}

func readAllMaps() [](map[string]int) {
	files, err := os.ReadDir("tmp")
	if err != nil {
		log.Fatal(err)
	}

	res := [](map[string]int){}

	for _, file := range files {
		k := fmt.Sprintf("tmp/%s", file.Name())
		res = append(res, decodeDict(k))
	}

	return res
}

func merge(maps [](map[string]int), out *Output) map[string]int {
	var wg sync.WaitGroup
	for _, m := range maps {
		wg.Add(1)
		go func(l map[string]int) {
			defer wg.Done()
			for k, v := range l {
				out.mu.Lock()
				out.out_map[k] += v
				out.mu.Unlock()
			}
			fmt.Printf("length of final map: %d\n", len(out.out_map))
		}(m)
	}
	wg.Wait()
	return out.out_map
}

func writeOut(output map[string]int) {
	jsonString, err := json.Marshal(output)
	check(err)

	_ = ioutil.WriteFile("output.json", jsonString, 0644)
}

func Reduce() {
	out := Output{
		out_map: map[string]int{},
	}
	slice_maps := readAllMaps()
	output := merge(slice_maps, &out)
	writeOut(output)
	fmt.Println(len(output))
}
