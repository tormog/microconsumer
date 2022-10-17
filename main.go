package main

import (
	"fmt"
	"log"
	"math/rand"
	tc "microconsumer/internal/consumers"
	tp "microconsumer/internal/producers"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

func envLoad() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local"
	}
	envFile := fmt.Sprintf("%s/.env.%s", dir, env)
	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Could not load .env.%s file: %v\n", env, err)
	}
}

func newRandomName(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	envLoad()

	numConsumers := 1
	if os.Getenv("NUMBER_OF_CONSUMERS") != "" {
		_consumers, err := strconv.Atoi(os.Getenv("NUMBER_OF_CONSUMERS"))
		if err != nil {
			log.Fatalln("Error reading NUMBER_OF_CONSUMERS. Please check syntax.")

		}
		numConsumers = _consumers
	}

	queueStoreName := os.Getenv("QUEUE_STORE_NAME")
	queueStoreDLName := os.Getenv("QUEUE_STORE_DL_NAME")

	log.Println("Starting producers and consumers")
	var wg sync.WaitGroup
	wg.Add(1)
	go tp.ProducerService(&wg, newRandomName(10))
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		go tc.ConsumerService(&wg, queueStoreName, queueStoreDLName, newRandomName(10), tc.SQLStore)
	}
	wg.Wait()
	log.Println("Main function done")
}
