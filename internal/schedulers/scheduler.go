package schedulers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lauralee01/orbit/internal/evaluator"
	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"sync"
	"time"
)

// FactsProvider is a pluggable function that supplies facts for a scheduled run.
// The scheduler does NOT know where facts come from — it delegates that responsibility.
// Returning nil means "skip this run" (e.g., missing data, external API down).
type FactsProvider func(ctx context.Context, ruleset storage.StoredRuleset) (rules.Facts, error)

// Scheduler manages all scheduled rule evaluations.
// It dynamically registers cron jobs based on rulesets stored in the database.
// It also updates or removes jobs when rulesets change.
type Scheduler struct {
	db            *sql.DB       // Database connection for loading rulesets + rules
	cron          *cron.Cron    // The cron engine that actually runs scheduled jobs
	factsProvider FactsProvider // Function that provides facts for evaluation
	dispatcher    evaluator.EvaluationDispatcher
	mu            sync.Mutex             // Protects entries + lastSpec from concurrent access
	entries       map[int64]cron.EntryID // Maps ruleset ID → cron job ID (so we can remove/update jobs)
	lastSpec      map[int64]string       // Tracks last cron spec used for each ruleset (detects changes)
	interval      time.Duration
}

// New creates a new Scheduler instance.
// It initializes the cron engine and internal maps.
func New(db *sql.DB, factsProvider FactsProvider, interval time.Duration) *Scheduler {
	return &Scheduler{
		db:            db,
		cron:          cron.New(), // Create a new cron scheduler
		factsProvider: factsProvider,
		dispatcher:    &evaluator.DirectDispatcher{DB: db},
		entries:       make(map[int64]cron.EntryID),
		lastSpec:      make(map[int64]string),
		interval:      interval,
	}
}

// Start launches the scheduler.
//  1. Starts the cron engine.
//  2. Immediately loads + registers all scheduled rulesets.
//  3. Starts a background loop that refreshes schedules every minute.
//     This allows dynamic updates without restarting the server.
func (s *Scheduler) Start(ctx context.Context) {
	s.cron.Start() // Start cron engine goroutines

	// Initial load of all schedules
	if err := s.refresh(ctx); err != nil {
		log.Printf("scheduler: initial refresh failed: %v", err)
	}

	// Periodically refresh schedules so changes in DB take effect automatically
	ticker := time.NewTicker(s.interval)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Server shutting down → stop scheduler cleanly
				s.Stop()
				return

			case <-ticker.C:
				// Refresh schedules every minute
				if err := s.refresh(ctx); err != nil {
					log.Printf("scheduler: refresh failed: %v", err)
				}
			}
		}
	}()
}

// Stop shuts down the cron engine.
// This stops all scheduled jobs from running.
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// refresh reloads all rulesets from the DB and updates cron jobs accordingly.
// It handles:
// - adding new jobs
// - updating changed jobs
// - removing disabled/deleted jobs
// This is the heart of dynamic scheduling.
func (s *Scheduler) refresh(ctx context.Context) error {
	rulesets, err := storage.ListRulesets(ctx, s.db)
	if err != nil {
		return fmt.Errorf("list rulesets: %w", err)
	}

	// Tracks which rulesets SHOULD have active cron jobs
	active := make(map[int64]struct{})

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, rs := range rulesets {

		// Skip rulesets without scheduling enabled
		if !rs.ScheduleEnabled {
			continue
		}

		// Validate cron expression + timezone and build final spec
		spec, err := validateAndBuildSpec(rs.ScheduleCron, rs.ScheduleTZ)
		if err != nil {
			log.Printf("scheduler: ruleset %d invalid schedule: %v", rs.ID, err)
			continue
		}

		active[rs.ID] = struct{}{}

		// If schedule hasn't changed, skip re-registering
		if prev, ok := s.lastSpec[rs.ID]; ok && prev == spec {
			continue
		}

		// Remove old cron job if it exists
		if id, ok := s.entries[rs.ID]; ok {
			s.cron.Remove(id)
			delete(s.entries, rs.ID)
		}

		// Capture rs for closure
		ruleset := rs

		// Register new cron job
		entryID, addErr := s.cron.AddFunc(spec, func() {
			s.runRuleset(context.Background(), ruleset)
		})
		if addErr != nil {
			log.Printf("scheduler: add job ruleset %d: %v", rs.ID, addErr)
			continue
		}

		// Save metadata for future refresh cycles
		s.entries[rs.ID] = entryID
		s.lastSpec[rs.ID] = spec

		log.Printf("scheduler: registered ruleset %d with spec %q", rs.ID, spec)
	}

	// Remove cron jobs for rulesets that are no longer active
	for rulesetID, entryID := range s.entries {
		if _, ok := active[rulesetID]; ok {
			continue // still active
		}

		// Remove job + metadata
		s.cron.Remove(entryID)
		delete(s.entries, rulesetID)
		delete(s.lastSpec, rulesetID)

		log.Printf("scheduler: removed ruleset %d job", rulesetID)
	}

	return nil
}

// validateAndBuildSpec ensures the cron expression + timezone are valid.
// It returns a full cron spec including CRON_TZ= prefix.
// Example output: "CRON_TZ=America/New_York 0 9 * * *"
func validateAndBuildSpec(cronExpr, tz string) (string, error) {
	cronExpr = strings.TrimSpace(cronExpr)
	tz = strings.TrimSpace(tz)

	if cronExpr == "" {
		return "", fmt.Errorf("empty cron expression")
	}
	if tz == "" {
		tz = "UTC" // default timezone
	}

	// Validate cron syntax
	if _, err := cron.ParseStandard(cronExpr); err != nil {
		return "", fmt.Errorf("invalid cron %q: %w", cronExpr, err)
	}

	// Validate timezone
	if _, err := time.LoadLocation(tz); err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", tz, err)
	}

	// Build final cron spec with timezone
	return "CRON_TZ=" + tz + " " + cronExpr, nil
}

// runRuleset executes a scheduled evaluation.
// Steps:
// 1. Load rules for the ruleset
// 2. Fetch facts via the FactsProvider
// 3. Evaluate rules
// 4. Log result
func (s *Scheduler) runRuleset(ctx context.Context, rs storage.StoredRuleset) {
	// Ensure we have a facts provider
	if s.factsProvider == nil {
		log.Printf("scheduler: ruleset %d skipped: no facts provider configured", rs.ID)
		return
	}

	// Fetch facts for this run
	facts, err := s.factsProvider(ctx, rs)
	if err != nil {
		log.Printf("scheduler: ruleset %d facts provider error: %v", rs.ID, err)
		return
	}
	if facts == nil {
		log.Printf("scheduler: ruleset %d skipped: facts provider returned nil facts", rs.ID)
		return
	}

	// Evaluate rules
	evalCtx := evaluator.EvalContext{
		TriggerSource: evaluator.TriggerSourceSchedule,
		TriggerAt:     time.Now(),
	}

	s.dispatcher.DispatchEvaluation(ctx, rs, facts, evalCtx)
}
