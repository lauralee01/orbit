package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lauralee01/orbit/internal/facts"
	"github.com/lauralee01/orbit/internal/handlers"
	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/schedulers"
	"github.com/lauralee01/orbit/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Printf("godotenv: %v (using environment variables only)", err)
    }

    // Open DB once per process
    db, err := storage.Open(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // --- Create the real facts provider ---
    factsProviderDB := facts.NewDBProvider(db)

    factsProvider := func(ctx context.Context, rs storage.StoredRuleset) (rules.Facts, error) {
        return factsProviderDB.GetFacts(ctx, rs.ID)
    }

    // --- Create cancellable root context ---
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // --- Handle OS signals for graceful shutdown ---
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // --- HTTP routes ---
    mux := http.NewServeMux()
    mux.HandleFunc("/health", health)
    mux.HandleFunc("GET /api/rulesets", handlers.ListRulesets(db))
    mux.HandleFunc("POST /api/rulesets", handlers.CreateRuleset(db))
    mux.HandleFunc("GET /api/rules", handlers.ListRules(db))
    mux.HandleFunc("POST /api/rules", handlers.CreateRule(db))
    mux.HandleFunc("POST /api/evaluate", handlers.Evaluate(db))

    // --- Scheduler ---
    scheduler := schedulers.New(db, factsProvider)
    scheduler.Start(ctx)

    // --- Server ---
    addr := ":8080"
    if p := os.Getenv("PORT"); p != "" {
        addr = ":" + p
    }

    server := &http.Server{
        Addr:    addr,
        Handler: mux,
    }

    go func() {
        log.Printf("listening on %s", addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("HTTP server: %v", err)
        }
    }()

    // --- Wait for shutdown signal ---
    <-sigCh
    log.Println("shutting down...")
    cancel()
   
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer shutdownCancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Printf("server shutdown error: %v", err)
    }

    log.Println("shutdown complete")

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
