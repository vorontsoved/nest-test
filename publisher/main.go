package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type server struct {
	nc *nats.Conn
}

func main() {
	var s server
	var err error

	uri := "nats://localhost:4222"

	for i := 0; i < 5; i++ {
		nc, err := nats.Connect(uri)
		if err == nil {
			s.nc = nc
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}

	js, err := s.nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "123",
		Subjects: []string{"123.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	js.Publish("123.TTTopic", []byte("TT"))

}
