package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

func grep(data string) string {
	start := time.Now()
	result := []string{}
	var path string = "./dbs/eksmo.csv"

	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(f), "\n")

	re := regexp.MustCompile(data)

	for _, line := range lines {
		if re.MatchString(line) {
			result = append(result, line)
		}
	}

	//fmt.Println(result)
	if len(result) > 0 {
		elapsed := time.Since(start)
		fmt.Printf("Time took %s\n", elapsed)
		return strings.Join(result, "</br>")
	} else {
		return "Not found"
	}
}
