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
		"https://uttc-hackathon-frontend-pink.vercel.app":                                true,
	}

	// back-onchainからのリクエストはOriginがないか、異なるOriginの可能性がある
	// その場合はすべてのOriginを許可（サーバー間通信のため）
	if origin == "" || allowedOrigins[origin] {
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Originがない場合は*を設定（サーバー間通信）
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
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

// loggingResponseWriter はレスポンスのステータスコードを記録するためのラッパー
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}
