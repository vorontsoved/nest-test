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

	// Установка соединения с сервером NATS
	nc, err := nats.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Создание JetStream контекста
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Создание контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Вызываем cancel() после завершения работы горутины

	js.AddStream(&nats.StreamConfig{
		Name:      "123",
		Subjects:  []string{"123.*"},
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Запуск горутины для асинхронного получения сообщений
	go func() {
		// Подписываемся на тему и получаем сообщения в реальном времени
		sub, err := js.Subscribe("123.TTTopic", func(msg *nats.Msg) {
			msg.Ack()
			m, err := msg.Metadata()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Это первый \n")
			fmt.Printf("Топик: %v\n", m.Stream)
			fmt.Printf("Группа: %v\n", msg.Subject)
			fmt.Printf("Получено новое сообщение: %v\n", string(msg.Data))
		}, nats.Durable("asss"), nats.DeliverNew())
		if err != nil {
			log.Fatal(err)
		}
		defer sub.Unsubscribe()

		// Ожидание сигнала завершения или отмены контекста
		<-ctx.Done()
		fmt.Println("Горутина получения сообщений завершает работу")
	}()

	// Ожидание сигнала завершения программы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Программа завершает работу")
}
