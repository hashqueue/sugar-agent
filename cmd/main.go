package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"sugar-agent/pkg/task"
	"sugar-agent/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

// doWork do the work
// messages: message channel
// return: none
func doWork(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Printf("[x] Received a message [x] -> %s", d.Body)
		// Unmarshal msg data
		msg := make(map[string]interface{})
		err := json.Unmarshal(d.Body, &msg)
		utils.FailOnError(err, "Failed to unmarshal message")
		baseUrl := msg["metadata"].(map[string]interface{})["base_url"].(string)
		taskId := msg["metadata"].(map[string]interface{})["task_id"].(string)

		// Login to get token
		loginData := map[string]interface{}{
			"username": msg["metadata"].(map[string]interface{})["username"],
			"password": msg["metadata"].(map[string]interface{})["password"],
		}
		token, err := utils.UserLogin(baseUrl, loginData)
		utils.FailOnError(err, "Failed to login")

		// update task status to RECEIVED
		updateData := map[string]interface{}{
			"task_status": 1, //RECEIVED
		}
		err = utils.UpdateTaskStatus(baseUrl, updateData, taskId, token)
		utils.FailOnError(err, "Failed to update task status")

		log.Printf("[x] Start task [x]")

		// update task status to STARTED
		updateData = map[string]interface{}{
			"task_status": 2, //STARTED
		}
		err = utils.UpdateTaskStatus(baseUrl, updateData, taskId, token)
		utils.FailOnError(err, "Failed to update task status")

		bT := time.Now()
		taskStatus := 3 //SUCCESS
		resultDesc := "everything is ok"
		data, err := task.StartTask(d.Body)
		if err != nil {
			taskStatus = 4 //FAILURE
			resultDesc = err.Error()
		}
		log.Printf("[x] Task is done [x]")
		log.Printf("[x] Total use time: %f s [x]", time.Since(bT).Seconds())

		// update task status to SUCCESS or FAILURE
		result := map[string]interface{}{
			"data": data,
			"desc": resultDesc,
		}
		updateData = map[string]interface{}{
			"task_status": taskStatus,
			"result":      result,
		}
		err = utils.UpdateTaskStatus(baseUrl, updateData, taskId, token)
		utils.FailOnError(err, "Failed to update task status")

		err = d.Ack(false)
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
	startConsuming("guest", "guest", "localhost", "5672", "device_exchange", "collect_device_perf_data_queue", "device_perf_data")

}
