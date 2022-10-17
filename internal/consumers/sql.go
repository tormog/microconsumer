package consumers

import (
	"encoding/json"
	"log"
	md "microconsumer/internal/models"
	dt "microconsumer/internal/storage"
)

type sqlStoreBackend struct {
	ConsumerData
}

func (sql *sqlStoreBackend) Consume() error {
	s, err := dt.NewSQLStorage()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("in consumer %s redis length:%v\n", sql.consumerID, sql.queue.Length(sql.queueStoreName))
	for sql.queue.Length(sql.queueStoreName) > 0 {
		result, err := sql.queue.Pop(sql.queueStoreName)
		if err != nil {
			return err
		}
		var redisData md.DataStore
		if err := json.Unmarshal([]byte(result), &redisData); err != nil {
			return err
		}
		if err := s.PersistData(redisData.ID, redisData.Type, redisData); err != nil {
			if err := sql.queue.Push(sql.queueDLStoreName, result); err != nil {
				log.Printf("consumer %s push to DL queue %s failed:%v\n", sql.consumerID, sql.queueDLStoreName, err)
			}
			log.Fatalf("consumer %s insert error:%v\n", sql.consumerID, err)
		}
		log.Printf("%s consumed %s id:%s\n", sql.consumerID, redisData.Type, redisData.ID)
	}
	return nil
}
