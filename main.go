package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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

func grep(data string, msgType int, conn *websocket.Conn) {
	var dbNames = listFiles()
	// resultData := []string{}
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
		conn.WriteMessage(msgType, []byte(match))
		// resultData = append(resultData, match)
		fmt.Println(match)
	}
	// return strings.Join(resultData, "</br></br>")
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
		conn.WriteMessage(msgType, []byte("Search started!"))
		grep(string(msg), msgType, conn)

		// Write message back to browser
		conn.WriteMessage(msgType, []byte("Search ended"))
		fmt.Println("----------------------------------------------------------------")
		return
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/submit.html")

	// outputText := map[string]string{"outputText": grep(r.Form.Get("data"))}
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/echo", handleEcho)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
