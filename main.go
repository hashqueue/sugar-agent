package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// failOnError check error
// err: error
// msg: error message
// return: none
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("Error happened, %s: %s", msg, err)
	}
}

// doWork do the work
// messages: message channel
// return: none
func doWork(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Printf("[x] Received a message [x] -> %s", d.Body)

		time.Sleep(5 * time.Second)

		log.Printf("[x] Task is done [x]")
		err := d.Ack(false)
		failOnError(err, "Failed to ack message")
	}
}

// startConsuming connect to MQ server
// user: MQ user name
// password: MQ user password
// host: MQ server host
// port: MQ server port
// exchange: MQ exchange name
// queue: MQ queue name
// return: none
func startConsuming(user string, password string, host string, port string, exchangeName string, queueName string, routingKey string) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		failOnError(err, "Failed to close connection")
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		failOnError(err, "Failed to close channel")
	}(ch)

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// set prefetchCount to 1: only one message will be sent to a worker at a time
	// set prefetchSize to 0: no effect
	// set global to false: the QoS settings apply to the current channel only
	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, exchangeName, routingKey)
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go doWork(messages)

	log.Printf("[******] Started consumer [******] -> Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	startConsuming("guest", "guest", "localhost", "5672", "logs_direct", "log_queue", "info")
}
