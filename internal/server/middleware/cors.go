package middleware

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
)

func CORS() func(connectHandler http.Handler) http.Handler {
	return func(connectHandler http.Handler) http.Handler {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: connectcors.AllowedMethods(),
			AllowedHeaders: append([]string{"Authorization"}, connectcors.AllowedHeaders()...),
			ExposedHeaders: connectcors.ExposedHeaders(),
			MaxAge:         7200, // 2 hours in seconds
		})
		return c.Handler(connectHandler)
	}
}
