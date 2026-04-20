package schedulers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lauralee01/orbit/internal/rules"
	"github.com/lauralee01/orbit/internal/storage"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"sync"
	"time"
)

// FactsProvider resolves the facts to evaluate for a scheduled ruleset run.
// Returning nil facts means "skip this run" and is logged by the scheduler.
type FactsProvider func(ctx context.Context, ruleset storage.StoredRuleset) (rules.Facts, error)

type Scheduler struct {
	db            *sql.DB
	cron          *cron.Cron
	factsProvider FactsProvider

	mu       sync.Mutex
	entries  map[int64]cron.EntryID
	lastSpec map[int64]string
}

func New(db *sql.DB, factsProvider FactsProvider) *Scheduler {
	return &Scheduler{
		db:            db,
		cron:          cron.New(),
		factsProvider: factsProvider,
		entries:       make(map[int64]cron.EntryID),
		lastSpec:      make(map[int64]string),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.cron.Start()
	if err := s.refresh(ctx); err != nil {
		log.Printf("scheduler: initial refresh failed: %v", err)
	}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				s.Stop()
				return
			case <-ticker.C:
				if err := s.refresh(ctx); err != nil {
					log.Printf("scheduler: refresh failed: %v", err)
				}
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) refresh(ctx context.Context) error {
	rulesets, err := storage.ListRulesets(ctx, s.db)
	if err != nil {
		return fmt.Errorf("list rulesets: %w", err)
	}

	active := make(map[int64]struct{})

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, rs := range rulesets {
		if !rs.ScheduleEnabled {
			continue
		}
		spec, err := validateAndBuildSpec(rs.ScheduleCron, rs.ScheduleTZ)
		if err != nil {
			log.Printf("scheduler: ruleset %d invalid schedule: %v", rs.ID, err)
			continue
		}

		active[rs.ID] = struct{}{}
		if prev, ok := s.lastSpec[rs.ID]; ok && prev == spec {
			continue
		}
		if id, ok := s.entries[rs.ID]; ok {
			s.cron.Remove(id)
			delete(s.entries, rs.ID)
		}

		ruleset := rs
		entryID, addErr := s.cron.AddFunc(spec, func() {
			s.runRuleset(context.Background(), ruleset)
		})
		if addErr != nil {
			log.Printf("scheduler: add job ruleset %d: %v", rs.ID, addErr)
			continue
		}

		s.entries[rs.ID] = entryID
		s.lastSpec[rs.ID] = spec
		log.Printf("scheduler: registered ruleset %d with spec %q", rs.ID, spec)
	}

	for rulesetID, entryID := range s.entries {
		if _, ok := active[rulesetID]; ok {
			continue
		}
		s.cron.Remove(entryID)
		delete(s.entries, rulesetID)
		delete(s.lastSpec, rulesetID)
		log.Printf("scheduler: removed ruleset %d job", rulesetID)
	}

	return nil
}

func validateAndBuildSpec(cronExpr, tz string) (string, error) {
	cronExpr = strings.TrimSpace(cronExpr)
	tz = strings.TrimSpace(tz)
	if cronExpr == "" {
		return "", fmt.Errorf("empty cron expression")
	}
	if tz == "" {
		tz = "UTC"
	}

	if _, err := cron.ParseStandard(cronExpr); err != nil {
		return "", fmt.Errorf("invalid cron %q: %w", cronExpr, err)
	}
	if _, err := time.LoadLocation(tz); err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", tz, err)
	}

	return "CRON_TZ=" + tz + " " + cronExpr, nil
}

func (s *Scheduler) runRuleset(ctx context.Context, rs storage.StoredRuleset) {
	storedRules, err := storage.ListRulesByRulesetID(ctx, s.db, rs.ID)
	if err != nil {
		log.Printf("scheduler: ruleset %d list rules: %v", rs.ID, err)
		return
	}

	ruleSlice := make(rules.Rules, len(storedRules))
	for i, row := range storedRules {
		ruleSlice[i] = rules.Rule{
			Field:    row.Field,
			Operator: row.Operator,
			Value:    row.Value,
		}
	}

	if s.factsProvider == nil {
		log.Printf("scheduler: ruleset %d skipped: no facts provider configured", rs.ID)
		return
	}
	facts, err := s.factsProvider(ctx, rs)
	if err != nil {
		log.Printf("scheduler: ruleset %d facts provider error: %v", rs.ID, err)
		return
	}
	if facts == nil {
		log.Printf("scheduler: ruleset %d skipped: facts provider returned nil facts", rs.ID)
		return
	}

	ok, evalErr := rules.Evaluate(facts, ruleSlice)
	if evalErr != nil {
		log.Printf("scheduler: ruleset %d evaluation error: %v", rs.ID, evalErr)
		return
	}
	log.Printf("scheduler: ruleset %d evaluation ok=%v", rs.ID, ok)
}
