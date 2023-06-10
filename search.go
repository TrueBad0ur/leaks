package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func readFile(filename string, cChannel chan string, data string) {
	defer wg.Done()
	//fmt.Println("Filename: ", filename)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	lines := strings.Split(string(body), "\n")
	re := regexp.MustCompile(data)
	for _, line := range lines {
		if re.MatchString(line) {
			cChannel <- line
		}
	}
}

func grep(data string) string {
	resultData := []string{}
	cChannel := make(chan string)
	filesCount := 5

	for i := 0; i < filesCount; i++ {
		wg.Add(1)
		go readFile(fmt.Sprintf("./tests/file%d", i), cChannel, data)
	}

	go func() {
		wg.Wait()
		close(cChannel)
	}()

	for match := range cChannel {
		resultData = append(resultData, match)
		fmt.Println(match)
	}
	return strings.Join(resultData, "</br>")
}
