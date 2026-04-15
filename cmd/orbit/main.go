package main

// Phase 4 (persistence + evaluate) is done. Phase 5 (hardening): see cmd/orbit/phase5_tasks.go and
// internal/handlers/phase5_tasks.go, internal/rules/phase5_tasks.go, internal/storage/phase5_tasks.go.

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lauralee01/orbit/internal/handlers"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("godotenv: %v (using environment variables only)", err)
	}

	// Open DB once per process; pass `db` into handlers per rulesets_tasks.go (closures or struct).
	db, err := storage.Open(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("GET /api/rulesets", handlers.ListRulesets(db))
	mux.HandleFunc("POST /api/rulesets", handlers.CreateRuleset(db))
	mux.HandleFunc("GET /api/rules", handlers.ListRules(db))
	mux.HandleFunc("POST /api/rules", handlers.CreateRule(db))
	mux.HandleFunc("POST /api/evaluate", handlers.Evaluate(db))

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}

}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}
