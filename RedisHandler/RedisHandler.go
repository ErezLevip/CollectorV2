package RedisHandler

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/bsm/redis-lock"
)

type RedisHandler struct {
	client *redis.Client
	config RedisConfiguration
}

type RedisConfiguration struct {
	Db               int
	Credentials      string
	ConnectionString string
}

type CacheHandler interface {
	Init(config RedisConfiguration)
	Lock(key string) (newLock *lock.Locker, err error)
	UnLock(lockedLock *lock.Locker)
	Get(key string) (string, error)
	Set(key string, val interface{}, duration time.Duration) error
	Delete(key string) error
}

func (self *RedisHandler) init() {
	self.client = redis.NewClient(&redis.Options{
		Addr:     self.config.ConnectionString,
		Password: self.config.Credentials,
		DB:       self.config.Db})
}
func (self *RedisHandler) Init(config RedisConfiguration) {
	self.config = config

}

func (self *RedisHandler) Lock(key string) (newLock *lock.Locker, err error) {
	self.init()
	defer self.client.Close()
	newLock, err = lock.Obtain(self.client, key, &lock.Options{
		LockTimeout: time.Duration(5) * time.Second,
		RetryDelay:   time.Duration(300) * time.Microsecond,
	})
	return
}

func (self *RedisHandler) UnLock(lockedLock *lock.Locker) {
	self.init()
	defer self.client.Close()
	if lockedLock.IsLocked() {
		lockedLock.Unlock()
	}
}

func (self *RedisHandler) Get(key string) (string, error) {
	self.init()
	defer self.client.Close()
	res := self.client.Get(key)
	if res.Err() != nil {
		return "", res.Err()
	}
	if res.Val() != "" {
		return res.Val(), nil
	}
	return "", nil
}

func (self *RedisHandler) Set(key string, val interface{}, duration time.Duration) error {
	self.init()
	defer self.client.Close()
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	status := self.client.Set(key, string(jsonBytes), duration)
	return status.Err()
}

func (self *RedisHandler) Delete(key string) error {
	self.init()
	status := self.client.Del(key)
	defer self.client.Close()
	return status.Err()
}
