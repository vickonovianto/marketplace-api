package main

import (
	"log"
	"marketplace-api/config"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}

	config := config.NewConfig()
	server := InitServer(config)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	wg.Wait()
}
