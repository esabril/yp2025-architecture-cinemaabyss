package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {
	rs := rand.New(rand.NewSource(time.Now().UnixNano()))

	port := getEnv("PORT", "8000")

	monolithUrl := getEnv("MONOLITH_URL", "http://127.0.0.1:8080")
	moviesServiceUrl := getEnv("MOVIES_SERVICE_URL", "http://127.0.0.1:8081")
	eventsServiceUrl := getEnv("EVENTS_SERVICE_URL", "http://127.0.0.1:8082")

	gradualMigration, _ := strconv.ParseBool(getEnv("GRADUAL_MIGRATION", "true"))
	moviesMigrationPercent, _ := strconv.ParseInt(getEnv("MOVIES_MIGRATION_PERCENT", "50"), 10, 64)

	handleHealth()
	handleEventsRequests(eventsServiceUrl)
	handleMoviesRequests(rs, gradualMigration, moviesMigrationPercent, monolithUrl, moviesServiceUrl)
	handleApiRequests(monolithUrl)

	log.Printf(
		"Starting proxy microservice on port %s\nMonolith Url: %s\nMovies Service Url: %s\nGradual Migration: %v\n",
		port, monolithUrl, moviesServiceUrl, gradualMigration,
	)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHealth() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"status": true})
	})
}

func handleMoviesRequests(rs *rand.Rand, gradualMigration bool, moviesMigrationPercent int64, monolithUrl, moviesServiceUrl string) {
	http.HandleFunc("/api/movies", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Movies Requests: graduationEnabled=%v, graduationPercent=%d\n", gradualMigration, moviesMigrationPercent)

		if useProxyToMovieService(rs, gradualMigration, moviesMigrationPercent) {
			log.Println("🎥 Proxy movies request to Movies Service:", r.Method, r.URL.Path)

			proxyToService(moviesServiceUrl, w, r)

			return
		}

		log.Println("⚙ Proxy movies request to Monolith:", r.Method, r.URL.Path)

		proxyToService(monolithUrl, w, r)
	})

	http.HandleFunc("/api/movies/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("🎥 Proxy movies request to Movies Service:", r.Method, r.URL.Path)

		proxyToService(moviesServiceUrl, w, r)

		return
	})
}

func handleEventsRequests(eventsServiceUrl string) {
	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		log.Println("✉ Proxy request to Events Service:", r.Method, r.URL.Path)

		proxyToService(eventsServiceUrl, w, r)
	})
}

func handleApiRequests(monolithUrl string) {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("⚙ Proxy request to Monolith:", r.Method, r.URL.Path)

		proxyToService(monolithUrl, w, r)
	})
}

func proxyToService(rawUrl string, w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ServeHTTP(w, r)
}

func useProxyToMovieService(rs *rand.Rand, gradualMigration bool, moviesMigrationPercent int64) bool {
	if !gradualMigration {
		// Если флаг выключен, маршрутизируем весь трафик в соответствующий микросервис
		return true
	}

	if moviesMigrationPercent > 100 {
		moviesMigrationPercent = 100
	}

	return rs.Intn(100) < int(moviesMigrationPercent)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
