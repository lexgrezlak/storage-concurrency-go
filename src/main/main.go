package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"storage-api/src/database"
	"storage-api/src/handler"
	"storage-api/src/service"
	"time"
	"github.com/gorilla/mux"
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

	router.HandleFunc("/promotions/{id}", handler.GetPromotionById(api)).Methods(http.MethodGet)

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
			redis.Set(context.Background(), p.Id, b,30*time.Minute)
		}
	}()


	srv := &http.Server{
		Addr:              "localhost:1321",
		Handler:           router,
		TLSConfig:         nil,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	log.Printf("Listening at: %v", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
