package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"sugar-agent/pkg/task"
	"sugar-agent/pkg/utils"
)

// doWork do the work
// messages: message channel
// return: none
func doWork(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Printf("[x] Received a message [x] -> %s", d.Body)
		bT := time.Now()
		task.StartTask(d.Body)
		log.Printf("[x] Task is done [x]")
		log.Printf("[x] Total use time: %f s [x]", time.Since(bT).Seconds())
		err := d.Ack(false)
		utils.FailOnError(err, "Failed to ack message")
	}
}

// startConsuming connect to MQ server
// user: MQ username
// password: MQ user password
// host: MQ server host
// port: MQ server port
// exchangeName: MQ exchange name
// queueName: MQ queue name
// routingKey: MQ routing key
// return: none
func startConsuming(user string, password string, host string, port string, exchangeName string, queueName string, routingKey string) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port))
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		utils.FailOnError(err, "Failed to close connection")
	}(conn)

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		utils.FailOnError(err, "Failed to close channel")
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
	utils.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	// set prefetchCount to 1: only one message will be sent to a worker at a time
	// set prefetchSize to 0: no effect
	// set global to false: the QoS settings apply to the current channel only
	err = ch.Qos(1, 0, false)
	utils.FailOnError(err, "Failed to set QoS")

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, exchangeName, routingKey)
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	utils.FailOnError(err, "Failed to bind a queue")

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go doWork(messages)

	log.Printf("[******] Started consumer [******] -> Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	startConsuming("guest", "guest", "localhost", "5672", "logs_direct", "log_queue", "info")
}
