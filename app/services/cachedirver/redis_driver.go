package cachedirver

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/rs/zerolog/log"
)

type RedisDriver struct {
	client *redis.Client
	mutex  *redsync.Mutex
}

func (driver *RedisDriver) SetConnection(conn interface{}) {
	driver.client = conn.(*redis.Client)
}

func (driver *RedisDriver) Set(key string, value interface{}, expiration time.Duration) error {
	switch value.(type) {
	case string:
		{
			st := driver.client.Set(context.Background(), key, value, expiration)
			if st.Err() != nil {
				log.Err(st.Err()).
					Str("func", "RedisDriver.Set").
					Str("@", "driver.client.Set").
					Send()
				return st.Err()
			}
			return nil
		}
	default:
		{
			bts, err := json.Marshal(&value)
			if err != nil {
				log.Err(err).
					Str("func", "RedisDriver.Set").
					Str("@", "json.Marshal").
					Send()
				return err
			}
			st := driver.client.Set(context.Background(), key, string(bts), expiration)
			if st.Err() != nil {
				log.Err(st.Err()).
					Str("func", "RedisDriver.Set").
					Str("@", "driver.client.Set").
					Send()
				return st.Err()
			}
			return nil
		}
	}
}

func (driver *RedisDriver) Get(key string) (str string, err error) {
	value := driver.client.Get(context.Background(), key)
	if value == nil {
		err := errors.New("ErrCacheRecordNotFound")
		return "", err
	}
	if value.Err() != nil {
		return "", value.Err()
	}
	return value.Val(), nil
}

func (driver *RedisDriver) UnmarshalToObject(key string, object interface{}) error {
	value, err := driver.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), object)
}

func (driver *RedisDriver) Delete(key string) error {
	res := driver.client.Del(context.Background(), key)
	if res.Err() != nil {
		log.Err(res.Err()).
			Str("func", "RedisDriver.Delete").
			Str("@", "driver.client.Del").
			Send()
		return res.Err()
	}
	d, err := res.Result()
	if err != nil {
		log.Err(err).
			Str("func", "RedisDriver.Delete").
			Str("@", "res.Result()").
			Send()
		return err
	}
	if d == 0 {
		err := errors.New("ErrCacheRecordNotFound")
		log.Err(err).
			Str("func", "RedisDriver.Delete").
			Str("@", "d == 0").
			Send()
		return err
	}
	return nil
}

func (driver *RedisDriver) DeleteByPattern(key string) (deletedCount int64, err error) {
	scan := driver.client.Scan(context.Background(), 0, key+"*", 0)
	res, _, err := scan.Result()
	if err != nil {
		log.Err(err).
			Str("func", "RedisDriver.DeleteByPattern").
			Str("@", "scan.Result").
			Send()
		return 0, err
	}
	driver.client.Del(context.Background(), res...)

	return int64(len(res)), nil
}

func (driver *RedisDriver) ResetDB() error {
	status := driver.client.FlushDB(context.Background())
	if status.Err() != nil {
		log.Err(status.Err()).
			Str("func", "RedisDriver.ResetDB").
			Str("@", "status.Err").
			Send()
	}
	return status.Err()
}

func (driver *RedisDriver) Lock() error {
	return driver.mutex.Lock()
}
func (driver *RedisDriver) Unlock() error {
	_, err := driver.mutex.Unlock()
	if err != nil {
		log.Err(err).
			Str("func", "RedisDriver.Unlock").
			Str("@", "driver.mutex.Unlock").
			Send()
	}
	return err
}
