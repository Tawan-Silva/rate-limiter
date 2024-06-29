package ratelimiter

import (
	"time"
)

// LimitData godoc
// @Summary Struct to store rate limiter data
// @Description Struct to store rate limiter data
type LimitData struct {
	Key           string `json:"key"`
	Seconds       int64  `json:"seconds"`
	BlockDuration int64  `json:"block_duration"`
	MaxRequests   int64  `json:"max_requests"`
	Id            string `json:"id"`
}

type LimitDataInput struct {
	Key           string `json:"key"`
	Seconds       int64  `json:"seconds"`
	BlockDuration int64  `json:"block_duration"`
	MaxRequests   int64  `json:"max_requests"`
}

type Store interface {
	Increment(key string, seconds int64) (int64, error)
	SaveInfoLimitData(key string, data LimitData) error
	GetInfoLimitData(key string) (LimitData, error)
	SetBlockDuration(key string, value int64, expiration time.Duration) error
	GetBlockDuration(key string) (int64, error)
	UpdateLimitData(key string, data LimitDataInput) error
	GetAllLimitData() ([]LimitData, error)
}

type RateLimiter struct {
	store Store
}

func NewRateLimiter(store Store) *RateLimiter {
	return &RateLimiter{
		store: store,
	}
}

func (r *RateLimiter) SetLimitData(key string, data LimitData) error {
	return r.store.SaveInfoLimitData(key, data)
}

func (r *RateLimiter) GetLimitData(key string) (LimitData, error) {
	return r.store.GetInfoLimitData(key)
}

func (r *RateLimiter) UpdateLimitData(idGenerated string, data LimitDataInput) error {
	return r.store.UpdateLimitData(idGenerated, data)
}

func (r *RateLimiter) GetLimitDataByIdGenerated(idGenerated string) (LimitData, error) {
	return r.store.GetInfoLimitData(idGenerated)
}

func (r *RateLimiter) GetAllLimitData() ([]LimitData, error) {
	return r.store.GetAllLimitData()
}

func (r *RateLimiter) Block(key string, blockDuration time.Duration) error {
	return r.store.SetBlockDuration("blocked:"+key, 1, blockDuration)
}

func (r *RateLimiter) IsBlocked(key string) (bool, error) {
	val, err := r.store.GetBlockDuration("blocked:" + key)
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *RateLimiter) Limit(key string, limit int64, duration int64, blockDuration int64) (bool, error) {
	if blocked, _ := r.IsBlocked(key); blocked {
		return true, nil
	}

	count, err := r.store.Increment(key, duration)
	if err != nil {
		return false, err
	}

	if count > limit {
		err := r.Block(key, time.Duration(blockDuration)*time.Second)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
