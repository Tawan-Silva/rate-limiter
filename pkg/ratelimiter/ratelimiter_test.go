package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	IncrementFunc         func(key string, seconds int64) (int64, error)
	SaveInfoLimitDataFunc func(key string, data LimitData) error
	GetInfoLimitDataFunc  func(key string) (LimitData, error)
	SetBlockDurationFunc  func(key string, value int64, expiration time.Duration) error
	GetBlockDurationFunc  func(key string) (int64, error)
	UpdateLimitDataFunc   func(key string, data LimitDataInput) error
	GetAllLimitDataFunc   func() ([]LimitData, error)
}

func (m *MockStore) Increment(key string, seconds int64) (int64, error) {
	return m.IncrementFunc(key, seconds)
}

func (m *MockStore) SaveInfoLimitData(key string, data LimitData) error {
	return m.SaveInfoLimitDataFunc(key, data)
}

func (m *MockStore) GetInfoLimitData(key string) (LimitData, error) {
	return m.GetInfoLimitDataFunc(key)
}

func (m *MockStore) SetBlockDuration(key string, value int64, expiration time.Duration) error {
	return m.SetBlockDurationFunc(key, value, expiration)
}

func (m *MockStore) GetBlockDuration(key string) (int64, error) {
	return m.GetBlockDurationFunc(key)
}

func (m *MockStore) UpdateLimitData(key string, data LimitDataInput) error {
	return m.UpdateLimitDataFunc(key, data)
}

func (m *MockStore) GetAllLimitData() ([]LimitData, error) {
	return m.GetAllLimitDataFunc()
}

// TestSetLimitData tests the SetLimitData function
func TestSetLimitData(t *testing.T) {
	store := &MockStore{
		SaveInfoLimitDataFunc: func(key string, data LimitData) error {
			assert.Equal(t, "testKey", key)
			assert.Equal(t, LimitData{}, data)
			return nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	err := rateLimiter.SetLimitData("testKey", LimitData{})
	assert.NoError(t, err)
}

// TestGetLimitData tests the GetLimitData function
func TestGetLimitData(t *testing.T) {
	store := &MockStore{
		GetInfoLimitDataFunc: func(key string) (LimitData, error) {
			assert.Equal(t, "testKey", key)
			return LimitData{}, nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	data, err := rateLimiter.GetLimitData("testKey")
	assert.NoError(t, err)
	assert.Equal(t, LimitData{}, data)
}

// TestUpdateLimitData tests the UpdateLimitData function
func TestUpdateLimitData(t *testing.T) {
	store := &MockStore{
		UpdateLimitDataFunc: func(key string, data LimitDataInput) error {
			assert.Equal(t, "testKey", key)
			assert.Equal(t, LimitDataInput{}, data)
			return nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	err := rateLimiter.UpdateLimitData("testKey", LimitDataInput{})
	assert.NoError(t, err)
}

// TestGetLimitDataByIdGenerated tests the GetLimitDataByIdGenerated function
func TestGetLimitDataByIdGenerated(t *testing.T) {
	store := &MockStore{
		GetInfoLimitDataFunc: func(key string) (LimitData, error) {
			assert.Equal(t, "testKey", key)
			return LimitData{}, nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	data, err := rateLimiter.GetLimitDataByIdGenerated("testKey")
	assert.NoError(t, err)
	assert.Equal(t, LimitData{}, data)
}

// TestGetAllLimitData tests the GetAllLimitData function
func TestGetAllLimitData(t *testing.T) {
	store := &MockStore{
		GetAllLimitDataFunc: func() ([]LimitData, error) {
			return []LimitData{}, nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	data, err := rateLimiter.GetAllLimitData()
	assert.NoError(t, err)
	assert.Equal(t, []LimitData{}, data)
}

// TestBlock tests the Block function
func TestBlock(t *testing.T) {
	store := &MockStore{
		SetBlockDurationFunc: func(key string, value int64, expiration time.Duration) error {
			assert.Equal(t, "blocked:testKey", key)
			assert.Equal(t, int64(1), value)
			assert.Equal(t, time.Second, expiration)
			return nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	err := rateLimiter.Block("testKey", time.Second)
	assert.NoError(t, err)
}

// TestIsBlocked tests the IsBlocked function
func TestIsBlocked(t *testing.T) {
	store := &MockStore{
		GetBlockDurationFunc: func(key string) (int64, error) {
			assert.Equal(t, "blocked:testKey", key)
			return 0, nil
		},
	}
	rateLimiter := NewRateLimiter(store)

	blocked, err := rateLimiter.IsBlocked("testKey")
	assert.NoError(t, err)
	assert.False(t, blocked)
}
