package cache

type Cache interface {
	Get(key string, payload any) error
	Set(key string, data any) error
}
