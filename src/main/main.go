package main

import (
	"log"
	"os"
	"storage-api/src/database"
	"storage-api/src/service"
)


func main() {
	// Set up database client
	client, err := database.NewRedisClient()
	log.Println(client)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	f, err := os.Open("data/csv/promotions.csv")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()
	promotions, err := service.ReadPromotionsFromCSV(f)
	if err != nil {
		log.Fatalf("failed to read csv file: %v", err)
	}

	p := promotions[0]
	log.Println(p)

}