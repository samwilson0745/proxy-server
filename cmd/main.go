package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"proxy-server/helper"
	"proxy-server/server"
)

func main() {
	port := flag.Int("port", 3000, "Port on which the caching proxy server will run")
	origin := flag.String("origin", "", "URL of the origin server")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache")

	flag.Parse()
	helper.GetRedisClient()

	if *clearCache {
		err := helper.ClearRedis()
		if err != nil {
			log.Fatalf("Error clearing cache: %v\n", err)
		}
		fmt.Println("Cache cleared successfully")
		os.Exit(0)
	}

	if *origin == "" {
		log.Fatal("Origin URL must be provided using the --origin flag")
	}

	fmt.Printf("Starting caching proxy server on port %d, forwarding requests to %s\n", *port, *origin)
	server.Start(*port, *origin)
}
