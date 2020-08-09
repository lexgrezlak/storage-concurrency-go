package service

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

type Promotion struct {
	Id string
	Price float64
	Date string
}

func ReadPromotionsFromCSV(b io.Reader) ([]Promotion, error) {
	r := csv.NewReader(b)
	var promotions []Promotion
	_, err := r.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		id := record[0]
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("failed to parse float: %v", err)
			continue
		}
		date := record[2]


		p := Promotion{id, price, date}
		promotions = append(promotions, p)
	}

	return promotions, nil
}