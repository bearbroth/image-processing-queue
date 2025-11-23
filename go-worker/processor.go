package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func simulateWork() {
	time.Sleep(2 * time.Second)
}

func extract(msg amqp.Delivery) (string, string, error) {
	var data map[string]any
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		return "", "", err
	}

	jobID := data["jobID"].(string)
	image := data["image"].(string)

	return jobID, image, nil
}

func processor(jobID string) {

	// scope : add image param during real processing.

	fmt.Printf("[jobID = %s] started\n", jobID)

	if err := redisSetStatus(jobID, "processing"); err != nil {
		fmt.Println("redis error:", err)
	}

	simulateWork()

	if err := redisSetStatus(jobID, "done"); err != nil {
		fmt.Println("redis error:", err)
	}

	fmt.Printf("[jobID = %s] finished\n", jobID)
}
