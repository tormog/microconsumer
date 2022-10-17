package consumers

import (
	ch "microconsumer/internal/queue"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type BackendType int64

const (
	JStore BackendType = iota
	SQLStore
)

type Consumer interface {
	Consume() error
}

type ConsumerData struct {
	queue            ch.RedisCache
	queueStoreName   string
	queueDLStoreName string
	consumerID       string
}

func ConsumerService(wg *sync.WaitGroup, queueStoreName, queueDLStoreName, consumerID string, backend BackendType) error {
	defer wg.Done()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	queue := ch.NewQueue()
	ticker := time.NewTicker(time.Second * time.Duration(10))

	consumerData := ConsumerData{
		queue:            queue,
		queueStoreName:   queueStoreName,
		queueDLStoreName: queueDLStoreName,
		consumerID:       consumerID,
	}

	for {
		select {
		case s := <-c:
			log.Printf("Received signal %s for Consumer %s, finishing...\n", s, consumerID)
			return nil
		case <-ticker.C:
			switch backend {
			case SQLStore:
				sqlStore := sqlStoreBackend{consumerData}
				if err := sqlStore.Consume(); err != nil {
					log.Fatalf("error processing data for consumer %s, %v", consumerID, err)
				}
			case JStore:
				jStore := jStoreBackend{consumerData}
				if err := jStore.Consume(); err != nil {
					log.Fatalf("error processing data for consumer %s, %v", consumerID, err)
				}
			}
		}
	}
}
