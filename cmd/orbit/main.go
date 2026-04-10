package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/lauralee01/orbit/internal/rules"

	"github.com/lauralee01/orbit/internal/handlers"
)

func main() {
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
		{Field: "age", Operator: "equals", Value: "21"},
	}
	
	ok, err := rules.Evaluate(facts, ruleset)
	fmt.Println(ok, err)

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
