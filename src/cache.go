package lru

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Set(key string, value interface{})
	Remove(key string)
	Len() int
	Clear()
}

type entry struct {
	key   string
	value interface{}
}
