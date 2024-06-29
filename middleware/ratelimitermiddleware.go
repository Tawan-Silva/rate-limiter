package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"ratelimiter/configs"
	"ratelimiter/pkg/ratelimiter"
	"sync"
	"time"
)

// LimitDataInput godoc
// @Summary Struct to store rate limiter data for Swagger documentation
// @Description Struct to store rate limiter data for Swagger documentation
type LimitDataInput struct {
	Seconds       int64 `json:"seconds"`
	BlockDuration int64 `json:"block_duration"`
	MaxRequests   int64 `json:"max_requests"`
}

type LimitData = ratelimiter.LimitData

type RateLimiterMiddleware struct {
	rateLimiter                *ratelimiter.RateLimiter
	defaultLimitByIp           int64
	defaultRequestLimitInSec   int64
	defaultRequestLimitByToken int64
	defaultBlockDuration       time.Duration
	mutexes                    sync.Map
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewRateLimiterMiddleware(rateLimiter *ratelimiter.RateLimiter) *RateLimiterMiddleware {
	defaultLimitRequestsIp := configs.GetLimitRequestsDefaultByIP()
	defaultRequestLimitInSec := configs.GetRequestLimitInSec()
	defaultBlockDuration := time.Duration(configs.GetBlockDuration()) * time.Second
	return &RateLimiterMiddleware{
		rateLimiter:                rateLimiter,
		defaultLimitByIp:           defaultLimitRequestsIp,
		defaultRequestLimitInSec:   defaultRequestLimitInSec,
		defaultRequestLimitByToken: configs.GetLimitRequestsByToken(),
		defaultBlockDuration:       defaultBlockDuration,
	}
}

func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.getKey(r)
		mutex := m.getMutex(key)
		mutex.Lock()
		defer mutex.Unlock()

		if m.isRequestLimited(key, w) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UpdateRateLimiter godoc
// @Summary Update rate limiter settings
// @Description update rate limiter settings for a specific key (ip or token)
// @Tags rate limiter
// @Accept  json
// @Produce  json
// @Param body body LimitDataInput true "Update rate limiter settings"
// @Success 200 {object} LimitDataInput "Successfully updated rate limiter settings"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /update-rate-limiter/ [put]
// @Security ApiKeyAuth
func (m *RateLimiterMiddleware) UpdateRateLimiter(writer http.ResponseWriter, request *http.Request) {
	key := m.getKey(request)
	limitData, err := m.rateLimiter.GetLimitData(key)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&limitData)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = m.rateLimiter.SetLimitData(key, limitData)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(limitData)
}

// GetAllRateLimiter godoc
// @Summary Get all rate limiter settings
// @Description get all rate limiter settings
// @Tags rate limiter
// @Accept  json
// @Produce  json
// @Success 200 {array} LimitData "Successfully retrieved all rate limiter settings"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /get-all-rate-limiter [get]
func (m *RateLimiterMiddleware) GetAllRateLimiter(writer http.ResponseWriter, request *http.Request) {
	limitData, err := m.rateLimiter.GetAllLimitData()
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(limitData)
}

func (m *RateLimiterMiddleware) getKey(r *http.Request) string {
	key := r.Header.Get("API_KEY")
	if key == "" {
		key = r.RemoteAddr
	}
	return key
}

func (m *RateLimiterMiddleware) getMutex(key string) *sync.Mutex {
	mutex, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

func (m *RateLimiterMiddleware) isRequestLimited(key string, w http.ResponseWriter) bool {
	limitData, err := m.rateLimiter.GetLimitData(key)
	if err != nil || limitData.Seconds == 0 {
		maxReq := m.defaultLimitByIp
		if m.isToken(key) {
			maxReq = m.defaultRequestLimitByToken
		}
		limitData = ratelimiter.LimitData{
			Key:           key,
			Seconds:       m.defaultRequestLimitInSec,
			MaxRequests:   maxReq,
			BlockDuration: int64(m.defaultBlockDuration.Seconds()),
			Id:            uuid.New().String(),
		}
		err = m.rateLimiter.SetLimitData(key, limitData)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return true
		}
	}

	limited, err := m.rateLimiter.Limit(key, limitData.MaxRequests, limitData.Seconds, limitData.BlockDuration)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}

	if limited {
		w.WriteHeader(http.StatusTooManyRequests)
		m.writeErrorResponse(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
		return true
	}

	return false
}

func (m *RateLimiterMiddleware) getTotalReqLimit(key string, w http.ResponseWriter) int64 {
	totalReqLimit := m.defaultLimitByIp
	if key != "" {
		totalReqLimit = configs.GetLimitRequestsByToken()
		if !m.isTokenValid(key, w) {
			return 0
		}
	}
	return totalReqLimit
}

func (m *RateLimiterMiddleware) isTokenValid(key string, w http.ResponseWriter) bool {
	token, err := jwt.ParseWithClaims(key, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.GetSecretKey()), nil
	})

	if err != nil || !token.Valid {
		m.writeErrorResponse(w, "Invalid token")
		return false
	}

	claims := token.Claims.(*jwt.StandardClaims)
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		m.writeErrorResponse(w, "Token expired")
		return false
	}

	return true
}

func (m *RateLimiterMiddleware) writeErrorResponse(w http.ResponseWriter, message string) {
	response := ErrorResponse{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(response)
}

func (m *RateLimiterMiddleware) isToken(key string) bool {
	_, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.GetSecretKey()), nil
	})

	return err == nil
}
