package main

// Phase 4 (persistence): follow the step-by-step tasks in internal/storage/doc.go before
// adding DB calls here (Open DB, ping, then wire handlers or a small smoke test).

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
	"github.com/lauralee01/orbit/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/echo", handlers.Echo)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	facts := rules.Facts{
		"age": 21,
	}
	
	ruleset := rules.Rules{
		{Field: "age", Operator: "equals", Value: "25"},
	}
	
	ok, err := rules.Evaluate(facts, ruleset)
	fmt.Println(ok, err)

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
