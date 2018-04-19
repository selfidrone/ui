package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/selfidrone/messages"
)

var processing = false
var nc stan.Conn
var natsServer = flag.String("nats", "nats://localhost:4222", "connection string for nats server")
var source = flag.String("source", "/", "sourcefolder for web")
var latestMutex sync.Mutex

func main() {
	flag.Parse()

	var err error
	clientID := fmt.Sprintf("server-%d", time.Now().UnixNano())
	nc, err = stan.Connect("test-cluster", clientID, stan.NatsURL(*natsServer))
	if err != nil {
		log.Fatal("Unable to connect to nats server: ", err)
	}

	sub, _ := nc.Subscribe(messages.MessageLiveImage, func(m *stan.Msg) {
		go processMessage(m)
	})

	defer sub.Unsubscribe()

	startServer()

	handleExit()
}

func handleExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func processMessage(m *stan.Msg) {
	latestMutex.Lock()
	defer func() {
		processing = false
		latestMutex.Unlock()
	}()

	filename := "./latest.jpg"

	di := messages.DroneImage{}
	di.DecodeMessage(m.Data)

	data, _ := base64.StdEncoding.DecodeString(di.Data)
	tp := http.DetectContentType(data)
	log.Println("Got image", tp)

	if tp == "image/jpeg" || tp == "image/png" {
		di.SaveDataToFile(filename)
	}
}

func startServer() {
	fs := http.FileServer(http.Dir(*source))
	http.Handle("/", fs)
	http.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		if !nc.NatsConn().IsConnected() {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/live", func(rw http.ResponseWriter, r *http.Request) {
		latestMutex.Lock()
		defer latestMutex.Unlock()

		d, err := ioutil.ReadFile("./latest.jpg")
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		io.Copy(rw, bytes.NewBuffer(d))
	})

	log.Println("Listening...")
	http.ListenAndServe(":4000", nil)
}
