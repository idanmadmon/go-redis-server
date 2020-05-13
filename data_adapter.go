package go_redis_server

import (
	"errors"
	"sync"
)

type (
	safeMap struct {
		data map[string]string
		*sync.Mutex
	}

	DB struct {
		safeMap
		Cfg Redis
	}
)

func (db *DB) get(key string) (string, error) {
	if db.data == nil {
		return "", errors.New("DB isn't initialized")
	}

	val, ok := db.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

func (db *DB) set(key string, val string) error {
	if db.data == nil {
		return errors.New("DB isn't initialized")
	}

	if db.Cfg.DisableOverride && db.isExist(key) {
		return errors.New("key already exists")
	}

	db.Lock()
	db.data[key] = val
	db.Unlock()
	return nil
}

func (db *DB) isExist(key string) bool {
	if db.data == nil {
		return false
	}

	_, ok := db.data[key]
	return ok
}