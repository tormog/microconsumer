package consumers

import (
	"encoding/json"
	md "example/internal/models"
	dt "example/internal/storage"
	"log"
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
	log.Printf("redis length:%v\n", jStore.cache.Length(jStore.redisStoreKeyName))
	mStore := map[string]dt.Storage{}
	defer closeJStoreStorage(mStore)
	for jStore.cache.Length(jStore.redisStoreKeyName) > 0 {
		result, err := jStore.cache.Pop(jStore.redisStoreKeyName)
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
		store.PersistData(redisData.ID, redisData.Type, redisData.Data)
		log.Printf("consumer %s consumed %s id:%s\n", jStore.consumerID, redisData.Type, redisData.ID)
	}
	return nil
}
