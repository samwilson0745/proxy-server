package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type CacheEntry struct {
	Response []byte
	Headers  http.Header
	Expiry   time.Time
}

var cache = struct {
	m sync.Map
}{}

func Start(port int, origin string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.Method + ":" + r.URL.Path

		// Check cache
		if entry, found := cache.m.Load(cacheKey); found {
			e := entry.(CacheEntry)
			if e.Expiry.After(time.Now()) {
				log.Println("Cache HIT for", cacheKey)
				for key, values := range e.Headers {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.Header().Set("X-Cache", "HIT")
				w.Write(e.Response)
				return
			}
			cache.m.Delete(cacheKey)
		}

		// Forward request to origin
		log.Println("Cache MISS for", cacheKey)
		resp, err := http.Get(origin + r.URL.Path)
		if err != nil {
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from origin", http.StatusInternalServerError)
			return
		}

		// Cache response
		headers := make(http.Header)
		for key, values := range resp.Header {
			for _, value := range values {
				headers.Add(key, value)
			}
		}
		cache.m.Store(cacheKey, CacheEntry{
			Response: body,
			Headers:  headers,
			Expiry:   time.Now().Add(5 * time.Minute), // Cache TTL
		})

		// Add headers and return response
		for key, values := range headers {
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
