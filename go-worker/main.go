package main

import (
	"fmt"
)

func logerror(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}

func main() {

	fmt.Println("Worker started, waiting for jobs....\n")

	url := "amqp://guest:guest@localhost:5672/"
	initQueue, err := initQueue(url)
	logerror(err, "Failed to connect to RabbitMQ")
	defer initQueue.Conn.Close()
	defer initQueue.Channel.Close()

	msgs, err := initQueue.Consume()
	logerror(err, "Failed to register a consumer.")

	initializeRedisClient()

	for msg := range msgs {

		jobID, image, err := extract(msg)
		logerror(err, "Failed to unmarshal message body.")

		fmt.Printf("Job ID : %s\n\n Image : %s\n\n", jobID, image)
		processor(jobID)

		value, err := redisGetStatus(jobID)
		logerror(err, "Failed to get value from redis.")
		fmt.Printf("\nStatus : %s\n", value)

		fmt.Printf("\n[job = %s] --> cycle complete\n\n", jobID)

	}

}
