package main

import (
	"net/http"
)

func handleCors(next http.Handler, w http.ResponseWriter, r *http.Request) {
	// 1. CORSヘッダーの設定
	w.Header().Set("Access-Control-Allow-Origin", "...")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	next.ServeHTTP(w, r)
}
func corsMiddleware(next http.Handler) http.Handler {
	wrappedFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCors(next, w, r)
	})
	return wrappedFunc
}
