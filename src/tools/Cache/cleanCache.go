package Cache

import (
	"os"
	"path/filepath"
	"selfWeb/src/configuration"
	"sync"
	"time"
)

type CacheData struct {
	Data      map[string][]byte
	Timestamp time.Time
}

var (
	cache      sync.Map
	cacheMutex sync.Mutex
	expiryTime = 1 * time.Hour
)

func AddToCache(key string, data map[string][]byte) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	configuration.Logger.Info("cache files in memory")
	cache.Store(key, CacheData{Data: data, Timestamp: time.Now()})
}

func GetFromCache(key string) (map[string][]byte, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cachedData, ok := cache.Load(key)
	if !ok {
		return nil, false
	}
	configuration.Logger.Info("cache files in memory")
	return cachedData.(CacheData).Data, true
}

func ClearExpiredCache() {
	now := time.Now()
	var keysToDelete []interface{}

	cache.Range(func(key, value interface{}) bool {
		cacheData := value.(CacheData)
		if now.Sub(cacheData.Timestamp) > expiryTime {
			keysToDelete = append(keysToDelete, key)
		}
		return true
	})
	for _, key := range keysToDelete {
		cache.Delete(key)
	}
}

func StartCacheCleaner(duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			ClearExpiredCache()
		}
	}()
}

func CacheFilesInMemory(filesDir string) (map[string][]byte, error) {
	filesData := make(map[string][]byte)
	err := filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(filesDir, path)
			if err != nil {
				return err
			}
			filesData[relPath] = fileContent
		}
		return nil
	})
	return filesData, err
}
