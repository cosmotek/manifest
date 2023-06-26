package state

type StorageProvider interface {
	Put(key string, val any) error
	Get(key string) (any, error)
}
