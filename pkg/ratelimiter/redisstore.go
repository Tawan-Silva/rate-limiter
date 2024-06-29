package ratelimiter

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	client *redis.Client
}

func (r *RedisStore) UpdateLimitData(key string, data LimitDataInput) error {
	oldLimitData, err := r.GetInfoLimitData(key)
	if err != nil {
		log.Printf("Failed to get data for id %s: %v", key, err)
		return err
	}

	if data.Seconds != 0 {
		oldLimitData.Seconds = data.Seconds
	}

	if data.BlockDuration != 0 {
		oldLimitData.BlockDuration = data.BlockDuration
	}

	jsonData, err := json.Marshal(oldLimitData)
	if err != nil {
		log.Printf("Failed to marshal data for id %s: %v", key, err)
		return err
	}

	err = r.client.Set("info::"+key, jsonData, 0).Err()
	if err != nil {
		log.Printf("Failed to set id %s: %v", key, err)
		return err
	}
	return nil
}

func (r *RedisStore) GetAllLimitData() ([]LimitData, error) {
	keys, err := r.client.Keys("info::*").Result()
	if err != nil {
		log.Printf("Failed to get all keys: %v", err)
		return nil, err
	}

	var allData []LimitData
	for _, key := range keys {
		key = strings.TrimPrefix(key, "info::")
		data, err := r.GetInfoLimitData(key)
		if err != nil {
			log.Printf("Failed to get data for key %s: %v", key, err)
			return nil, err
		}
		allData = append(allData, data)
	}
	return allData, nil
}
func (r *RedisStore) SaveInfoLimitData(key string, data LimitData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data for key %s: %v", key, err)
		return err
	}

	err = r.client.Set("info::"+key, jsonData, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
		return err
	}
	return nil
}

func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) Ping() error {
	return r.client.Ping().Err()
}

func (r *RedisStore) Increment(key string, seconds int64) (int64, error) {
	val, err := r.client.IncrBy("limit::"+key, 1).Result()
	if err != nil {
		log.Printf("Failed to increment key %s: %v", key, err)
		return 0, err
	}

	go r.client.Expire("limit::"+key, time.Duration(seconds)*time.Second).Result()
	return val, nil
}

func (r *RedisStore) SetBlockDuration(key string, value int64, expiration time.Duration) error {
	err := r.client.Set(key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set block duration for key %s: %v", key, err)
		return err
	}

	return nil
}

func (r *RedisStore) GetBlockDuration(key string) (int64, error) {
	val, err := r.client.Get(key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		log.Printf("Failed to get block duration for key %s: %v", key, err)
		return 0, err
	}

	return val, nil
}

func (r *RedisStore) SaveLimitData(key string, data LimitData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data for key %s: %v", key, err)
		return err
	}

	err = r.client.Set(key, jsonData, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
		return err
	}
	return nil
}

func (r *RedisStore) GetInfoLimitData(key string) (LimitData, error) {
	val, err := r.client.Get("info::" + key).Result()
	if err != nil {
		log.Printf("Failed to get key %s: %v", key, err)
		return LimitData{}, err
	}

	var data LimitData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Printf("Failed to unmarshal data for key %s: %v", key, err)
		return LimitData{}, err
	}
	return data, nil
}
