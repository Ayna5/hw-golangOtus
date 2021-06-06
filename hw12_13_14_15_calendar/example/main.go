package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var configFile string

type QueueTask struct {
	Number1 int
	Number2 int
}

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewSender(configFile)
	if err != nil {
		log.Fatalf("can't get config1: %v", err)
	}

	conn, err := amqp.Dial(config.MQ.URI)
	handleError(err, "cannot connect to AMQP")
	fmt.Println("connected to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "cannot create a amqpChannel")
	fmt.Println("created a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("queue", true, false, false, false, nil)
	handleError(err, "could not declare `queue` queue")
	fmt.Println("declare queue:", queue)

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "could not configure QoS")
	fmt.Println("configured QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "could not register consumer")
	fmt.Println("consumer is register")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			addTask := &QueueTask{}

			err := json.Unmarshal(d.Body, addTask)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("Result of %d + %d is : %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// Остановка для завершения программы
	<-stopChan
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
