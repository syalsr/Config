package main

import (
	"Config/app/config"
	"Config/app/database"
	"Config/app/server"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.GetConfig()
	err := database.NewClient(context.TODO(), cfg)
	log.Println("Connected to PostgreSQL")
	if err != nil {
		log.Println(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go server.StartGrpcServer(&wg)
	go server.StartHppServe(&wg)

	wg.Wait()
}

func init() {
	err := godotenv.Load("/app/.env")
	fmt.Println(*config.GetConfig())
	if err != nil {
		log.Fatal(err)
	}
}
