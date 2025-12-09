package main

import (
	"net/http"
)

func mainHandler(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-frontend-ix9mgv4j6-kazukis-projects-f47db0d9.vercel.app/sell")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
