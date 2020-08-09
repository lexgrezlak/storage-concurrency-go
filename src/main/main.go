package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"storage-api/src/config"
	"storage-api/src/database"
	"storage-api/src/handler"
	"storage-api/src/service"
	"time"
)

const (
	REDIS_EXPIRATION = 30 * time.Minute
)

func main() {
	// Initialize config. If it can't find the file, it will load the variables
	// from the environment. It would be a good idea to read the file path to the config
	// from environment, because we might want to have `test.yml` or some other config.
	c, err := config.GetConfig("development.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("config has been loaded: %v", c)

	// Set up Redis.
	redis, err := database.NewRedisClient(c.Redis)
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

	// For the purposes of this problem we'll read the CSV just from the file saved on disk,
	// and the filename will be a hardcoded value. Therefore you need to run the program
	// from the project's root directory, that is `go run src/main/main.go`.
	// Open the csv file and read its records.
	f, err := os.Open("promotions.csv")
	if err != nil {
		log.Fatalf("failed to open csv file: %v", err)
	}
	defer f.Close()
	promotions := service.GetPromotionsFromCSV(context.Background(), f)
	log.Printf("parsed and processed csv in %2f seconds", time.Since(start).Seconds())

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
