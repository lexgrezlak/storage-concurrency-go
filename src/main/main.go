package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"storage-api/src/database"
	"storage-api/src/handler"
	"storage-api/src/service"
	"time"
)

const (
	REDIS_EXPIRATION = 30 * time.Minute
)

func main() {
	// Set up Redis.
	redis, err := database.NewRedisClient()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Set up API.
	api := service.NewAPI(redis)

	// Set up router.
	router := mux.NewRouter()

	// Set up handlers.
	router.HandleFunc("/promotions/{id}", handler.GetPromotionById(api)).Methods(http.MethodGet)

	// Measure the performance of parsing and processing the CSV.
	start := time.Now()

	// Open the csv file and read its records.
	f, err := os.Open("promotions.csv")
	if err != nil {
		log.Fatalf("failed to open csv file: %v", err)
	}
	defer f.Close()
	promotions := service.GetPromotionsFromCSV(context.Background(), f)
	log.Printf("Parsed and processed csv in %2f seconds", time.Since(start).Seconds())

	// Set the records in Redis
	go func() {
		for p := range promotions {
			value, err := json.Marshal(p)
			if err != nil {
				log.Printf("failed to marshal promotion: %v", err)
			} else {
				key := p.Id
				redis.Set(context.Background(), key, value, REDIS_EXPIRATION)
			}
		}
	}()

	// Set up the server.
	srv := &http.Server{
		Addr:         "0.0.0.0:1321",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Listening at: %v", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
