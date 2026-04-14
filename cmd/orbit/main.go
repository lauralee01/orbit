package main

// Finish phase 4: implement persistence + HTTP per internal/storage/rules_crud_tasks.go,
// internal/handlers/rules_http_tasks.go, internal/handlers/evaluate_http_tasks.go — then register routes here.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/lauralee01/orbit/internal/handlers"
	"github.com/lauralee01/orbit/internal/storage"
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
	mux.HandleFunc("/api/echo", handlers.Echo)
	mux.HandleFunc("GET /api/rulesets", handlers.ListRulesets(db))
	mux.HandleFunc("POST /api/rulesets", handlers.CreateRuleset(db))

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
