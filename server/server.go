package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"proxy-server/helper"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheEntry struct {
	Response []byte
	Headers  http.Header
}

func Start(port int, origin string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.Method + ":" + r.URL.Path
		ctx := context.Background()
		client := helper.GetRedisClient()
		cached, err := client.Get(ctx, cacheKey).Result()
		if err == nil {
			var entry CacheEntry
			if err := json.Unmarshal([]byte(cached), &entry); err != nil {
				log.Printf("Error unmarshaling cache for %s : %v", cacheKey, err)
				client.Del(ctx, cacheKey)
			} else {
				log.Println("Cache hit for", cacheKey)
				for key, values := range entry.Headers {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.Header().Set("X-Cache", "HIT")
				w.Write(entry.Response)
				return
			}
		} else if err != redis.Nil {
			log.Printf("Redis error for %s: %v", cacheKey, err)

		}

		// Cache Miss - forward request to origin
		log.Println("Cache MISS for", cacheKey)

		resp, err := http.Get(origin + r.URL.Path)

		if err != nil {
			log.Println("Error while getting response")
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from origin", http.StatusInternalServerError)
		}

		entry := CacheEntry{
			Response: body,
			Headers:  resp.Header,
		}

		jsonData, err := json.Marshal(entry)

		if err != nil {
			log.Printf("Error marshalling cache entry for %s: %v", cacheKey, err)
		} else {
			err = client.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()

			if err != nil {
				log.Printf("Error storing in Redis for %s: %v", cacheKey, err)
			}
		}

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.Header().Set("X-Cache", "MISS")
		w.Write(body)
	})
	log.Printf("Proxy server running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
