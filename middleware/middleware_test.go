package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"
	"ratelimiter/configs"
	"ratelimiter/pkg/ratelimiter"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	configs.LoadConfig()
	store := ratelimiter.NewRedisStore("localhost:6379")
	rateLimiter := ratelimiter.NewRateLimiter(store)
	middleware := NewRateLimiterMiddleware(rateLimiter)

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	t.Run("allows correct number of requests", func(t *testing.T) {
		for i := 0; i < int(configs.GetLimitRequestsDefaultByIP()); i++ {
			req := httptest.NewRequest("GET", "/home", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		}
	})

	t.Run("blocks after limit reached", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/home", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusTooManyRequests {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
		}
	})

	t.Run("unblocks after block duration", func(t *testing.T) {
		time.Sleep(time.Duration(configs.GetBlockDuration()) * time.Second)
		req := httptest.NewRequest("GET", "/home", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("allows correct number of requests by token", func(t *testing.T) {
		expirationTime := time.Duration(configs.GetExpirationToken()) * time.Second
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		})

		tokenString, err := token.SignedString([]byte(configs.GetSecretKey()))
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		for i := 0; i < int(configs.GetLimitRequestsByToken()); i++ {
			req := httptest.NewRequest("GET", "/home", nil)

			req.Header.Add("API_KEY", tokenString)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		}
	})

	t.Run("blocks after limit reached by token", func(t *testing.T) {
		expirationTime := time.Duration(configs.GetExpirationToken()) * time.Second
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		})

		tokenString, err := token.SignedString([]byte(configs.GetSecretKey()))
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/home", nil)
		req.Header.Add("API_KEY", tokenString)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusTooManyRequests {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
		}
	})

	t.Run("unblocks after block duration by token", func(t *testing.T) {
		time.Sleep(time.Duration(configs.GetBlockDuration()) * time.Second)
		expirationTime := time.Duration(configs.GetExpirationToken()) * time.Second
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		})

		tokenString, err := token.SignedString([]byte(configs.GetSecretKey()))
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest("GET", "/home", nil)
		req.Header.Add("API_KEY", tokenString)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
