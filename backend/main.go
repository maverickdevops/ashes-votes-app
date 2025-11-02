package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
)

// Struct for main application
type App struct {
	db *sql.DB
}

// Struct representing incoming vote request
type VoteRequest struct {
	Team string `json:"team"`
}

// Struct representing team counts for the /counts endpoint
type Count struct {
	Team  string `json:"team"`
	Count int64  `json:"count"`
}

// Allowed teams
var allowedTeams = map[string]bool{"australia": true, "england": true}

func main() {
	// --- Initialize OpenTelemetry ---
	ctx := context.Background()
	cleanup, err := InitOTel(ctx)
	if err != nil {
		log.Fatalf("failed to init otel: %v", err)
	}
	defer cleanup(ctx) // ensures flush of traces on shutdown

	// --- Environment variables for DB and port ---
	pgURL := getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/votes?sslmode=disable")
	port := getenv("PORT", "8080")

	// --- Connect to Postgres ---
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}

	// --- Wait until DB is reachable ---
	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := db.Ping(); err != nil {
			if time.Now().After(deadline) {
				log.Fatalf("db ping failed: %v", err)
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	app := &App{db: db}

	// --- Register HTTP endpoints ---
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/vote", app.voteHandler)
	http.HandleFunc("/counts", app.countsHandler)

	// --- Start the HTTP server ---
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, allowCORS(http.DefaultServeMux)))
}

// getenv returns env variable or default
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// Health check endpoint for container readiness
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// --- Vote handler ---
// Accepts POST /vote { "team": "australia" }
// Inserts a vote into the DB and records a trace span
func (a *App) voteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("vote-handler").Start(r.Context(), "cast-vote")
	defer span.End() // close span when function ends

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var vr VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&vr); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Validate team name
	if !allowedTeams[vr.Team] {
		http.Error(w, "invalid team", http.StatusBadRequest)
		return
	}

	// Insert vote into Postgres
	_, err := a.db.ExecContext(ctx, "INSERT INTO votes (team) VALUES ($1)", vr.Team)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// --- Counts handler ---
// Returns JSON array of current vote counts for both teams
func (a *App) countsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query("SELECT team, COUNT(*) FROM votes GROUP BY team")
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// default 0 values
	counts := map[string]int64{"australia": 0, "england": 0}
	for rows.Next() {
		var team string
		var cnt int64
		rows.Scan(&team, &cnt)
		counts[team] = cnt
	}

	resp := []Count{
		{"australia", counts["australia"]},
		{"england", counts["england"]},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- allowCORS ---
// Allows frontend JS app to call backend APIs from another port (localhost:3000 → 8080)
func allowCORS(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	}
}
