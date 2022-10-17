package queue

type Queue interface {
	Push(key string, value interface{}) error
	Pop(key string) (string, error)
	Length(key string) int64
}
