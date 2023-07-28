package cachedirver

import (
	"time"

	"github.com/stretchr/testify/mock"
)

func InitCacheDriverMock() *CacheDriverMock {
	m := CacheDriverMock{}
	conn = &m
	return &m
}

type CacheDriverMock struct {
	mock.Mock
}

func (m *CacheDriverMock) SetConnection(conn interface{}) {}

func (m *CacheDriverMock) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *CacheDriverMock) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *CacheDriverMock) GetByPattern(pattern string) (map[string]string, error) {
	args := m.Called(pattern)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *CacheDriverMock) UnmarshalToObject(key string, object interface{}) error {
	args := m.Called(key, object)
	return args.Error(0)
}

func (m *CacheDriverMock) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *CacheDriverMock) DeleteByPattern(key string) (deletedCount int64, err error) {
	args := m.Called(key)
	return int64(args.Int(0)), args.Error(1)
}

func (m *CacheDriverMock) Lock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *CacheDriverMock) Unlock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *CacheDriverMock) ResetDB() error {
	args := m.Called()
	return args.Error(0)
}
