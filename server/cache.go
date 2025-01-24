package server

import (
	"sync"
)

func ClearCache() error {
	cache.m = sync.Map{} // Reset the cache
	return nil
}
