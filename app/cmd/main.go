package main

import (
	"Config/app/config"
	"Config/app/database"
	"Config/app/server"
	"context"
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

	database.Migrate(cfg)
	log.Printf("Databse migration successfully")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go server.StartGrpcServer(&wg)
	go server.StartHppServe(&wg)

	wg.Wait()
}

func init() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf(".env loaded")
}
