package consumers

import (
	"encoding/json"
	"log"
	md "microconsumer/internal/models"
	dt "microconsumer/internal/storage"
)

type jStoreBackend struct {
	ConsumerData
}

func closeJStoreStorage(mStores map[string]dt.Storage) {
	for k := range mStores {
		if err := mStores[k].Close(); err != nil {
			log.Printf("issues saving storage backend, %v\n", err)
		}
	}
}

func (jStore *jStoreBackend) Consume() error {
	log.Printf("redis length:%v\n", jStore.queue.Length(jStore.queueStoreName))
	mStore := map[string]dt.Storage{}
	defer closeJStoreStorage(mStore)
	for jStore.queue.Length(jStore.queueStoreName) > 0 {
		result, err := jStore.queue.Pop(jStore.queueStoreName)
		if err != nil {
			return err
		}
		var redisData md.DataStore
		if err := json.Unmarshal([]byte(result), &redisData); err != nil {
			return err
		}
		if _, ok := mStore[redisData.Type]; !ok {
			log.Printf("consumer %s loading datastore %s", jStore.consumerID, redisData.Type+".json.gz")
			mStore[redisData.Type] = dt.NewJSONStorage(redisData.Type + ".json.gz")
		}
		store := mStore[redisData.Type]
		if err := store.PersistData(redisData.ID, redisData.Type, redisData.Data); err != nil {
			if err := jStore.queue.Push(jStore.queueDLStoreName, result); err != nil {
				log.Printf("consumer %s push to DL queue %s failed:%v\n", jStore.consumerID, jStore.queueDLStoreName, err)
			}
			log.Fatalf("consumer %s insert error:%v\n", jStore.consumerID, err)
		}
		log.Printf("consumer %s consumed %s id:%s\n", jStore.consumerID, redisData.Type, redisData.ID)
	}
	return nil
}
