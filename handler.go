package main

import (
	"net/http"
)

func handleCors(next http.Handler, w http.ResponseWriter, r *http.Request) {
	// 1. CORSヘッダーの設定
	origin := r.Header.Get("Origin")
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true,
		"https://uttc-hackathon-frontend-ix9mgv4j6-kazukis-projects-f47db0d9.vercel.app": true,
	}

	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		// 実際の処理はせずに、CORSヘッダーを返して終了
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		return                              // 処理をここで終了させる
	}
	next.ServeHTTP(w, r)
}
func corsMiddleware(next http.Handler) http.Handler {
	wrappedFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCors(next, w, r)
	})
	return wrappedFunc
}
