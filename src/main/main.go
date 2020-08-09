package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"storage-api/src/database"
	"storage-api/src/service"
	"time"
)

func main() {
	// Set up database client
	client, err := database.NewRedisClient()
	log.Println(client, err)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Measure the performance of parsing and processing the CSV.
	start := time.Now()
	// Open the csv file and read its records.
	f, err := os.Open("data/csv/promotions.csv")
	if err != nil {
		log.Fatalf("failed to open csv file: %v", err)
	}
	defer f.Close()
	promotions, _ := service.GetPromotionsFromCSV(context.Background(), f)
	fmt.Printf("\n%2fs", time.Since(start).Seconds())

	start = time.Now()
	go func() {
		for p := range promotions {
			b, _ := json.Marshal(p)
			client.Set(context.Background(), p.Id, b,30*time.Minute)
		}
	}()

	//time.Sleep(30*time.Second)
	fmt.Printf("\n%2fs", time.Since(start).Seconds())
	for p := range promotions {
		str, err := client.Get(context.Background(), p.Id).Result()
		log.Printf("err: %v str %v", err, str)
	}


}
