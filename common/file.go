package common

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func ReadDataFromFile(filepath string) []string {
	f, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Errorf("failed to open file %s, error: %s\n", filepath, err.Error()))
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(fmt.Errorf("failed to close file %s, error: %s\n", filepath, err.Error()))
		}
	}(f)

	log.Printf("data is being loaded from path: %s, please wait\n", filepath)

	var lines []string
	count := 0
	rd := bufio.NewReader(f)
	for {
		privKey, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		lines = append(lines, strings.TrimSpace(privKey))
		count++
	}

	log.Printf("%d records are loaded\n", count)
	return lines
}
