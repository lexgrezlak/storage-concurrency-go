package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"strconv"
)

type Promotion struct {
	Id    string  `json:"id"`
	Price float64 `json:"price"`
	Date  string  `json:"date"`
}

// First pulls a promotion from redis which is a json string,
// and then transforms it into a Promotion interface,
// and returns it.
func (api *api) GetPromotionById(id string) (*Promotion, error) {
	jsonStr, err := api.r.Get(context.Background(), id).Result()
	if err != nil {
		return nil, err
	}

	var p Promotion
	if err = json.Unmarshal([]byte(jsonStr), &p); err != nil {
		log.Println(err)
		return nil, err
	}

	return &p, nil
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

func GetPromotionsFromCSV(ctx context.Context, b io.Reader) <-chan *Promotion {
	r := csv.NewReader(b)
	promotions := make(chan *Promotion)

	go func() {
		defer close(promotions)
		for {
			select {
			case <-ctx.Done():
				break
			default:
				record, err := r.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					// We'll skip the record that record that returned an error.
					log.Printf("failed to read record from csv: %v", err)
				} else {
					p, err := processRecordToPromotion(record)
					if err != nil {
						// We'll skip the record that record that returned an error.
						log.Printf("failed to process record to promotion: %v", err)
					} else {
						promotions <- p
					}
				}

			}
		}
	}()

	return promotions
}
