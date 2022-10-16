package consumers

import (
	ch "example/internal/cache"

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
	cache             ch.RedisCache
	redisStoreKeyName string
	consumerID        string
}

func ConsumerService(wg *sync.WaitGroup, redisStoreKeyName, consumerID string, backend BackendType) error {
	defer wg.Done()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	cache := ch.NewCache()
	ticker := time.NewTicker(time.Second * time.Duration(10))

	consumerData := ConsumerData{
		cache:             cache,
		redisStoreKeyName: redisStoreKeyName,
		consumerID:        consumerID,
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
