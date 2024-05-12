package main

import (
	"fmt"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMericsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits)))
}

func (cfg *apiConfig) resetZero(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func main() {
	port := "8050"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	cfg := apiConfig{
		fileserverHits: 0,
	}
	fileHander := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", cfg.middlewareMericsInc(fileHander))
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/metrics", cfg.getHits)
	mux.HandleFunc("/reset", cfg.resetZero)
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
