package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	uri := "nats://localhost:4222"

	nc, err := nats.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js.AddStream(&nats.StreamConfig{
		Name:     "123",
		Subjects: []string{"123.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		sub, err := js.QueueSubscribe("123.TTTopic", "Cons", func(msg *nats.Msg) {
			m, err := msg.Metadata()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Это первый \n")
			fmt.Printf("Топик: %v\n", m.Stream)
			fmt.Printf("Группа: %v\n", msg.Subject)
			fmt.Printf("Получено новое сообщение: %v\n", string(msg.Data))
		}, nats.DeliverNew())
		if err != nil {
			log.Fatal(err)
		}
		defer sub.Unsubscribe()

		<-ctx.Done()
		fmt.Println("Горутина получения сообщений завершает работу")
	}()

	go func() {
		sub, err := js.QueueSubscribe("123.TTTopic", "Cons", func(msg *nats.Msg) {
			m, err := msg.Metadata()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Это второй \n")
			fmt.Printf("Топик: %v\n", m.Stream)
			fmt.Printf("Группа: %v\n", msg.Subject)
			fmt.Printf("Получено новое сообщение: %v\n", string(msg.Data))
		}, nats.DeliverNew())
		if err != nil {
			log.Fatal(err)
		}
		defer sub.Unsubscribe()

		<-ctx.Done()
		fmt.Println("Горутина получения сообщений завершает работу")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Программа завершает работу")
}
