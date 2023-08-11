package cachedirver

import "time"

type CacheDriver interface {
	SetConnection(conn interface{})
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	UnmarshalToObject(key string, object interface{}) error
	Delete(key string) error
	DeleteByPattern(key string) (deletedCount int64, err error)
	Lock() error
	Unlock() error
	ResetDB() error
}
