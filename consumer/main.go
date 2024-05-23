package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	uri = "nats://nats:4222"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	consumerGroup := flag.String("group", "group_1", "consumer group")
	consumerName := flag.String("name", "group_name_1", "consumer name")
	topic := flag.String("topic", "topic_1", "topic to subscribe to")
	serverURI := flag.String("uri", uri, "server uri, possible separated by comma to indicate multiple servers")

	flag.Parse()
	logrus.Infof("group=%s, name=%s, topic=%s, uri=%q", *consumerGroup, *consumerName, *topic, *serverURI)

	cli, err := NewJetStreamClient(*serverURI, *consumerGroup, *consumerName)
	if err != nil {
		logrus.WithError(err).Fatal("error creating jetstream client")
	}

	type PublishMessage struct {
		Count int `json:"count,omitempty"`
	}

	expectedMessages := 3000
	result := make([]uint16, expectedMessages)
	resultLock := sync.Mutex{}
	cli.Subscribe(*topic, func(msg *Message) error {
		fmt.Printf("message: %v", msg.Value)
		pmsg := &PublishMessage{}
		if err := json.Unmarshal(msg.Value, pmsg); err != nil {
			logrus.WithError(err).Error("error unmarshaling message into PublishMessage")
			return nil // do not return error so that we won't receive this one again
		}

		resultLock.Lock()
		defer resultLock.Unlock()
		result[pmsg.Count]++
		return nil
	})

	ticker := time.NewTicker(time.Millisecond * 250)
	tickerFunc := func() (shouldStop bool) {
		resultLock.Lock()
		defer resultLock.Unlock()

		for i, count := range result {
			if count == 0 {
				logrus.Warnf("haven't yet received message %d", i)
				return false
			}
		}
		return true
	}
	for range ticker.C {
		if shouldStop := tickerFunc(); shouldStop {
			break
		}
	}

	// now range over the entire result slice to count if all messages have been received exactly once
	for i, count := range result {
		if count != 1 {
			logrus.Warnf("message %d has been received %d times", i, count)
		}
	}
}
