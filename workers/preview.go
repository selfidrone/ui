package workers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/selfidrone/messages"
)

type PreviewWorker struct {
	nc          stan.Conn
	sub         stan.Subscription
	latestMutex sync.Mutex
	imageFolder string
}

func NewPreviewWorker(natsConn, output string) *PreviewWorker {
	pw := &PreviewWorker{imageFolder: output}

	clientID := fmt.Sprintf("server-%d", time.Now().UnixNano())

	var err error
	pw.nc, err = stan.Connect("test-cluster", clientID, stan.NatsURL(natsConn))
	if err != nil {
		log.Fatal("Unable to connect to nats server: ", err)
	}

	return pw
}

func (p *PreviewWorker) Start() {
	log.Println("Starting preview worker")
	p.sub, _ = p.nc.Subscribe(messages.MessageDroneImage, p.processMessage)
}

func (p *PreviewWorker) Stop() {
	p.sub.Unsubscribe()
	p.sub.Close()
}

func (p *PreviewWorker) processMessage(m *stan.Msg) {
	p.latestMutex.Lock()
	defer func() {
		p.latestMutex.Unlock()
	}()

	filename := path.Join(p.imageFolder, "latest.jpg")

	di := messages.DroneImage{}
	di.DecodeMessage(m.Data)

	data, _ := base64.StdEncoding.DecodeString(di.Data)
	tp := http.DetectContentType(data)
	log.Println("Got image", tp)

	if tp == "image/jpeg" || tp == "image/png" {
		log.Println("Got image", filename)
		di.SaveDataToFile(filename)
	}
}
