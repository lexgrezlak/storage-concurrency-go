package service

import "github.com/go-redis/redis/v8"

type api struct {
	r *redis.Client
}

type PromotionDatastore interface {
	GetPromotionById(id string) (*Promotion, error)
}

func NewAPI(r *redis.Client) *api {
	return &api{r: r}
}