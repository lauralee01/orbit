package schedulers

import (
	"context"
	"database/sql"
	"github.com/lauralee01/orbit/internal/storage"
	"github.com/lauralee01/orbit/internal/handlers/evaluate"
	"github.com/robfig/cron/v3"
	"time"
	"log"
)

type Scheduler struct {
	db *sql.DB
	cron *cron.Cron
}

func StartScheduler(db *sql.DB {
	go func() {
		// load rulesets with schedule enabled
		rulesets, err := storage.ListRulesets(context.Background(), db)
		if err != nil {
			log.Printf("scheduler: list rulesets: %v", err)
		}

		for _, ruleset := range rulesets {
			if ruleset.ScheduleEnabled {
				scheduleCron := ruleset.ScheduleCron
				scheduleTZ := ruleset.ScheduleTZ
				schedule, err := cron.ParseStandard(scheduleCron)
				if err != nil {
					log.Printf("scheduler: parse cron: %v", err)
				}
				nextTime := schedule.Next(time.Now())
				  // run evaluations for this ruleset
				  evalCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				  defer cancel()
				  evalOK, evalReason := evaluate.Evaluate(evalCtx, db, ruleset.ID, nil)
				  if !evalOK {
					log.Printf("scheduler: evaluation failed: %s", evalReason)
				  } else {
					log.Printf("scheduler: evaluation succeeded")
				  }
			}
		}
	}
})