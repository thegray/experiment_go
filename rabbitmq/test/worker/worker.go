package main

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var (
	ExchangeName    = "aaaa"
	ExchangeType    = "direct"
	ExchangeDurable = true
	RoutingKey      = "zzzz"
	MandatoryFlag   = true
	ImmediateFlag   = true
	Q_ADDRESS       = "amqp://guest:guest@localhost:5672/"
	QueName         = "que1"
)

type Conn struct {
	Channel *amqp.Channel
}

func setupQueConn(rabbitURL string) (Conn, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to Connect to RabbitMQ at: %s", rabbitURL)
		return Conn{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to create channel")
		return Conn{}, err
	}

	_, err = ch.QueueDeclare(
		QueName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// err = ch.ExchangeDeclare(
	// 	ExchangeName,    // name
	// 	ExchangeType,    // type
	// 	ExchangeDurable, // durable
	// 	false,           // auto-deleted
	// 	false,           // internal
	// 	false,           // no-wait
	// 	nil,             // arguments
	// )
	// if err != nil {
	// 	log.Printf("Failed to declare an exchange with name: %s", ExchangeName)
	// 	return Conn{}, err
	// }

	return Conn{
		Channel: ch,
	}, err
}

func main() {
	conn, err := setupQueConn(Q_ADDRESS)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Channel.Close()

	msgs, err := conn.Channel.Consume(
		QueName, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(3 * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
