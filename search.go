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

func findInFile(filename string, cChannel chan string, data string) {
	defer wg.Done()

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

func listFiles() []string {
	var filesArray = []string{}

	files, err := ioutil.ReadDir("./dbs/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filesArray = append(filesArray, file.Name())
	}

	return filesArray
}

func grep(data string) string {
	var dbNames = listFiles()
	resultData := []string{}
	cChannel := make(chan string)

	for _, dbName := range dbNames {
		wg.Add(1)
		go findInFile("./dbs/"+dbName, cChannel, data)
	}

	go func() {
		wg.Wait()
		close(cChannel)
	}()

	for match := range cChannel {
		resultData = append(resultData, match)
		fmt.Println(match)
	}
	return strings.Join(resultData, "</br></br>")
}
