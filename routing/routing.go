package routing

import (
	"context"
	"fmt"
	"go-plate/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.Use(LangMiddleWare)
	router.Use(RequestLoggerMiddleware)
	User(CreateSub("/user", router))
}

func CreateSub(path string, router *mux.Router) *mux.Router {
	return router.PathPrefix(path).Subrouter()
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		id, err := services.ValidateToken(token)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), services.UserIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LangMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}

		ctx := context.WithValue(r.Context(), services.LangKey, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		services.Logger.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}