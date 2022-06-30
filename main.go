package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func reader(conn *websocket.Conn, msg StatusFinal1) {
	for {
		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// fmt.Println(string(p), "<<<")
		// fmt.Println(messageType, "<<<")
		if err := conn.WriteJSON(msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Successfully Connected")

	// uptimeTicker := time.NewTicker(5 * time.Second)
	// dateTicker := time.NewTicker(10 * time.Second)

	// for {
	// 	select {
	// 	case <-uptimeTicker.C:
	// 		reader(ws, run())
	// 	case <-dateTicker.C:
	// 		dateTicker.Stop()
	// 		return
	// 	}
	// }
	reader(ws, run())
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/status", wsEndpoint)
}

func main() {
	fmt.Println("udah jalan")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type StatusFinal struct {
	Water string `json:"water"`
	Wind  string `json:"wind"`
}

type StatusFinal1 struct {
	StatusFinal StatusFinal `json:"status"`
}

type Status1 struct {
	Status Status `json:"status"`
}

func run() StatusFinal1 {
	filePath := "./status.json"
	uptimeTicker := time.NewTicker(5 * time.Second)
	dateTicker := time.NewTicker(10 * time.Second)
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		select {
		case <-uptimeTicker.C:
			status := Status1{
				Status: Status{
					Water: randInt(1, 100),
					Wind:  randInt(1, 100),
				},
			}

			file, _ := json.MarshalIndent(status, "", " ")

			_ = ioutil.WriteFile("status.json", file, 0644)

			file, err1 := ioutil.ReadFile(filePath)
			if err1 != nil {
				fmt.Printf("// error while reading file %s\n", filePath)
				fmt.Printf("File error: %v\n", err1)
				os.Exit(1)
			}

			err2 := json.Unmarshal(file, &status)
			if err2 != nil {
				fmt.Println("error:", err2)
				os.Exit(1)
			}
			win := randInt(1, 100)
			wat := randInt(1, 100)
			statusFinal := StatusFinal1{
				StatusFinal: StatusFinal{
					Water: strconv.Itoa(win) + "m - " + wind(win),
					Wind:  strconv.Itoa(wat) + "s -" + water(wat),
				},
			}
			return statusFinal
		case <-dateTicker.C:
			uptimeTicker.Stop()
			dateTicker.Stop()
			return StatusFinal1{}
		}
	}

}

func water(val int) string {
	if val < 5 {
		return "aman"
	}
	if val >= 6 && val <= 8 {
		return "siaga"
	}
	if val > 8 {
		return "bahaya"
	}
	return ""
}

func wind(val int) string {
	if val < 6 {
		return "aman"
	}
	if val >= 7 && val <= 15 {
		return "siaga"
	}
	if val > 15 {
		return "bahaya"
	}
	return ""
}
