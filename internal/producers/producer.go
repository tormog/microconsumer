package producers

import (
	ch "example/internal/cache"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Producer interface {
	Produce(fromID string) (string, error)
}

type ProducerData struct {
	cache             ch.RedisCache
	redisStoreKeyName string
	resourceName      string
	producerID        string
}

func ProducerService(wg *sync.WaitGroup, producerID string) error {
	defer wg.Done()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(time.Second * time.Duration(10))

	tws := &twitterSource{
		producerData: ProducerData{
			cache:             ch.NewCache(),
			redisStoreKeyName: os.Getenv("REDIS_STORE_KEY_NAME"),
			resourceName:      os.Getenv("TWITTER_RESOURCE_NAME"),
			producerID:        producerID,
		},
		searchURL:   os.Getenv("TWITTER_SEARCH_URL"),
		searchQuery: os.Getenv("TWITTER_SEARCH_QUERY"),
		bearerToken: os.Getenv("TWITTER_BEARER_TOKEN"),
	}

	fromID := ""
	for {
		select {
		case s := <-c:
			log.Printf("Received signal %s for ProducerService %s, finishing...", s, producerID)
			return nil
		case <-ticker.C:
			newestID, err := tws.Produce(fromID)
			if err != nil {
				log.Fatalf("error processing data in twitter producer %s, %v", producerID, err)
			}
			fromID = newestID
		}
	}
}
