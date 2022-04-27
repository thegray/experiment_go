package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// server config
var (
	ExchangeName    = "aaaa"
	ExchangeType    = "direct"
	ExchangeDurable = true
	RoutingKey      = "zzzz"
	MandatoryFlag   = true
	ImmediateFlag   = true
	Q_ADDRESS       = "amqp://guest:guest@localhost:5672/"
	HTTP_PORT       = ":8000"
)

type Conn struct {
	Channel *amqp.Channel
}

type Service struct {
	QueCon *Conn
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

	err = ch.ExchangeDeclare(
		ExchangeName,    // name
		ExchangeType,    // type
		ExchangeDurable, // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Printf("Failed to declare an exchange with name: %s", ExchangeName)
		return Conn{}, err
	}

	return Conn{
		Channel: ch,
	}, err
}

func (conn Conn) publish(data string) error {
	return conn.Channel.Publish(
		ExchangeName, // exchange name
		RoutingKey,
		MandatoryFlag, // mandatory - we don't care if there is no queue
		ImmediateFlag, // immediate - we don't care if there is no consumer on the queue
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
			// DeliveryMode: amqp.Persistent,
		})
}

// func queSetup() {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	failOnError(err, "Failed to get Rabbit Channel")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"task_queue", // name
// 		true,         // durable
// 		false,        // delete when unused
// 		false,        // exclusive
// 		false,        // no-wait
// 		nil,          // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")
// }

func main() {
	conn, err := setupQueConn(Q_ADDRESS)
	failOnError(err, "Failed to connect to RabbitMQ")

	svc := Service{QueCon: &conn}

	// create a new echo instance
	e := echo.New()

	//Post Request
	e.POST("/api/experiment/postmessage", svc.PostMessage)

	e.Logger.Fatal(e.Start(HTTP_PORT))
}

func (svc Service) PostMessage(c echo.Context) error {
	type Message struct {
		Msg string `json:"msg"`
	}
	msg := Message{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&msg)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	log.Printf("received message %#v", msg)
	svc.QueCon.publish(msg.Msg)
	log.Printf("message queued")
	return c.JSON(http.StatusOK, Message{Msg: "we got your message: " + msg.Msg})
	// return c.String(http.StatusOK)
}
