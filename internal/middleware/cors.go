package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			//  we will set these later using the env.
			"http://127.0.0.1:5500",
			"http://localhost:5500",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(next)
}
