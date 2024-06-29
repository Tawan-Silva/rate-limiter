package server

import (
	"encoding/json"
	"github.com/pkg/browser"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	httpSwagger "github.com/swaggo/http-swagger"
	"ratelimiter/configs"
	_ "ratelimiter/docs" // Importa os documentos gerados pelo swag
	"ratelimiter/middleware"
)

type Server struct {
	rateLimiterMiddleware *middleware.RateLimiterMiddleware
}

type AuthTokenResponse struct {
	Token string `json:"token"`
}

type IndexResponse struct {
	Message string `json:"message"`
}

func NewServer(rateLimiterMiddleware *middleware.RateLimiterMiddleware) *Server {
	return &Server{
		rateLimiterMiddleware: rateLimiterMiddleware,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/token", s.GetAuthToken)
	mux.Handle("/swagger-ui/", httpSwagger.WrapHandler)
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	mux.Handle("/home", s.rateLimiterMiddleware.Middleware(http.HandlerFunc(s.Index)))
	//update-rate-limiter/${id}
	mux.HandleFunc("/update-rate-limiter/", s.rateLimiterMiddleware.UpdateRateLimiter)
	mux.HandleFunc("/get-all-rate-limiter", s.rateLimiterMiddleware.GetAllRateLimiter)
	// atualizar dados do rate limiter do ip ou token
	log.Println("Starting server on :8080")

	browser.OpenURL("http://localhost:8080/swagger-ui/")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

// Index godoc
// @Summary Welcome to the rate limited index page!
// @Description get index
// @Tags home
// @Accept  json
// @Produce  json
// @Success 200 {object} IndexResponse "Welcome to the rate limited index page!"
// @Router /home [get]
// @Security ApiKeyAuth
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	response := IndexResponse{
		Message: "Welcome to the rate limited index page!",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAuthToken godoc
// @Summary Generates a new auth token
// @Description get authToken
// @Tags token
// @Accept  json
// @Produce  json
// @Header 200 {string} Token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzI3NjU0MzF9.-I6C4q8FvuHqSvkYOhc5Yaoz2pPR_Use1z3hZJ4ETaE"
// @Success 200 {object} AuthTokenResponse "token detail"
// @Router /token [get]
func (s *Server) GetAuthToken(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Duration(configs.GetExpirationToken()) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.GetSecretKey()))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthTokenResponse{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
