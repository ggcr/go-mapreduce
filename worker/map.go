package worker

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(key string) string {
	v, err := os.ReadFile(key)
	check(err)
	fmt.Printf("file %s read ✅\n", key)
	return strip(string(v))
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		} else {
			result.WriteByte(' ')
		}
	}
	return result.String()
}

func genMap(file string, key string) map[string]int {
	tmpMap := make(map[string]int)
	for _, words := range strings.Fields(file) {
		tmpMap[words] += 1
	}
	return tmpMap
}

func writeDisk(tmpMap map[string]int, key string) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(tmpMap)
	check(err)

	f, err := os.Create(fmt.Sprintf("tmp/%s", filepath.Base(key)))
	check(err)

	defer f.Close()

	_, err = f.Write(b.Bytes())
	check(err)

	fmt.Printf("file %s wrote 📕\n", fmt.Sprintf("tmp/%s", filepath.Base(key)))
}

func Map(key string, wg *sync.WaitGroup) {
	defer wg.Done()
	v := readFile(key)
	fileMap := genMap(v, key)
	writeDisk(fileMap, key)
	fmt.Printf("%s DONE\n", key)
}
