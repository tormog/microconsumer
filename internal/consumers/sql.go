package consumers

import (
	"encoding/json"
	md "example/internal/models"
	dt "example/internal/storage"
	"log"
)

type sqlStoreBackend struct {
	ConsumerData
}

func (mysql *sqlStoreBackend) Consume() error {
	s, err := dt.NewSQLStorage()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("in consumer %s redis length:%v\n", mysql.consumerID, mysql.cache.Length(mysql.redisStoreKeyName))
	for mysql.cache.Length(mysql.redisStoreKeyName) > 0 {
		result, err := mysql.cache.Pop(mysql.redisStoreKeyName)
		if err != nil {
			return err
		}
		var redisData md.DataStore
		if err := json.Unmarshal([]byte(result), &redisData); err != nil {
			return err
		}
		err = s.PersistData(redisData.ID, redisData.Type, redisData)
		if err != nil {
			log.Fatalf("consumer %s insert error:%v\n", mysql.consumerID, err)
		}
		log.Printf("%s consumed %s id:%s\n", mysql.consumerID, redisData.Type, redisData.ID)
	}
	return nil
}
