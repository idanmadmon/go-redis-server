package go_redis_server

import (
	"errors"
	"sync"
)

type safeDB struct {
	data map[string]string
	*sync.Mutex
}

var DB safeDB

func get(key string) (string, error) {
	if DB.data == nil {
		return "", errors.New("DB isn't initialized")
	}

	val, ok := DB.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

func set(key, val string) error {
	if DB.data == nil {
		return errors.New("DB isn't initialized")
	}

	if cfg.DisableOverride && isExist(key) {
		return errors.New("key already exists")
	}

	DB.Lock()
	DB.data[key] = val
	DB.Unlock()
	return nil
}

func isExist(key string) bool {
	if DB.data == nil {
		return false
	}

	_, ok := DB.data[key]
	return ok
}