package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var wg sync.WaitGroup

func findInFile(filename string, cChannel chan string, data string) {
	defer wg.Done()

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile("(?i)" + data)

	for i := 1; scanner.Scan(); i++ {
		if re.MatchString(scanner.Text()) {
			cChannel <- scanner.Text()
			//fmt.Println(scanner.Text())
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

	// In this folder we have .gitkeep, which will always be first in the list
	filesArray = filesArray[1:]
	fmt.Println(filesArray)

	return filesArray
}

func grep(data string, msgType int, conn *websocket.Conn) {
	var dbNames = listFiles()
	for i := 0; i < len(dbNames); i++ {
		dbNames[i] = "./dbs/" + dbNames[i]
	}

	cChannel := make(chan string)

	for _, dbName := range dbNames {
		wg.Add(1)
		go findInFile(dbName, cChannel, data)
	}

	go func() {
		wg.Wait()
		close(cChannel)
	}()

	for match := range cChannel {
		conn.WriteMessage(msgType, []byte(match))
	}
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
		conn.WriteMessage(msgType, []byte("Search started!"))
		grep(string(msg), msgType, conn)

		conn.WriteMessage(msgType, []byte("Search ended"))
		fmt.Println("----------------------------------------------------------------")
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/submit.html")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/echo", handleEcho)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
