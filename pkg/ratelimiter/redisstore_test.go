package ratelimiter

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestUpdateLimitDataRedis tests the UpdateLimitData function
func TestUpdateLimitDataRedis(t *testing.T) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
	store := NewRedisStore(redisAddress)

	// Prepare the data
	err := store.SaveInfoLimitData("testKey", LimitData{})
	assert.NoError(t, err)

	// Test UpdateLimitData
	err = store.UpdateLimitData("testKey", LimitDataInput{})
	assert.NoError(t, err)

	// Clean up
	err = store.client.Del("info::testKey").Err()
	assert.NoError(t, err)
}

// TestGetAllLimitDataRedis tests the GetAllLimitData function
func TestGetAllLimitDataRedis(t *testing.T) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
	store := NewRedisStore(redisAddress)

	err := store.SaveInfoLimitData("testKey", LimitData{
		Seconds:       5,
		BlockDuration: 30,
		MaxRequests:   2,
		Id:            "testId",
	})
	assert.NoError(t, err)

	data, err := store.GetAllLimitData()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = store.client.Del("info::testKey").Err()
	assert.NoError(t, err)
}

// TestSaveInfoLimitDataRedis tests the SaveInfoLimitData function
func TestSaveInfoLimitDataRedis(t *testing.T) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
	store := NewRedisStore(redisAddress)

	err := store.SaveInfoLimitData("testKey", LimitData{
		Seconds:       5,
		BlockDuration: 30,
		MaxRequests:   2,
		Id:            "testId",
	})
	assert.NoError(t, err)

	data, err := store.GetInfoLimitData("testKey")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = store.client.Del("info::testKey").Err()
	assert.NoError(t, err)
}
