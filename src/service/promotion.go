package service

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

type Promotion struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Date string `json:"date"`
}

func processRecordToPromotion(record []string) (*Promotion, error) {
	id := record[0]

	price, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Printf("failed to parse float: %v", err)
		return nil, err
	}

	date := record[2]

	return &Promotion{id, price, date}, nil
}


func GetPromotionsFromCSV(ctx context.Context, b io.Reader) (<-chan *Promotion, <-chan error) {
	r := csv.NewReader(b)
	concurrency := 1000
	promotions := make(chan *Promotion, concurrency)
	errs := make(chan error, concurrency)

	go func() {
		defer close(promotions)
		defer close(errs)
		for {
			select {
			case <-ctx.Done():
				break
			default:
				record, err := r.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					errs <- err
				}

				p, err := processRecordToPromotion(record)
				if err != nil {
					errs <- err
				}

				promotions <- p
			}
		}
	}()

	return promotions, errs
}