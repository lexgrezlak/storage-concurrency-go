package main

import (
	"log"
	"os"
	"storage-api/src/service"
)

func main() {
	f, err := os.Open("data/promotions.csv")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	promotions, err := service.ReadPromotionsFromCSV(f)
	if err != nil {
		log.Fatalf("failed to read csv file: %v", err)
	}

	log.Println(promotions[0])
}