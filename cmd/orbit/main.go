package main

// Phase 4 (persistence): follow the step-by-step tasks in internal/storage/doc.go before
// adding DB calls here (Open DB, ping, then wire handlers or a small smoke test).
//
// (4) go.mod hygiene: when you add or remove imports, run `go mod tidy` from the repo root
// so direct dependencies are listed under `require` without spurious `// indirect` lines.

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

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/echo", handlers.Echo)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	// open database connection
	db, err := storage.Open(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create ruleset
	id, err := storage.CreateRuleset(context.Background(), db, "Test Ruleset")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Ruleset created with ID:", id)

	// list rulesets
	rulesets, err := storage.ListRulesets(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Rulesets:", rulesets)

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
