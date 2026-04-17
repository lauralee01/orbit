package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateRuleset_MinimalBody(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO rulesets`).
		WithArgs("policy", "", "", "UTC", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	h := CreateRuleset(db)
	body := `{"name":"policy"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rulesets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
	var out createRulesetResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.ID != 1 || out.Name != "policy" {
		t.Fatalf("response = %+v", out)
	}
}

func TestCreateRuleset_WithSchedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO rulesets`).
		WithArgs("policy", "", "0 9 * * *", "America/New_York", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	h := CreateRuleset(db)
	body := `{"name":"policy","schedule_cron":"0 9 * * *","schedule_tz":"America/New_York","schedule_enabled":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/rulesets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateRuleset_ValidationErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantSubstr string
	}{
		{
			name:       "empty name",
			body:       `{"name":""}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "name is required",
		},
		{
			name:       "enabled without cron",
			body:       `{"name":"x","schedule_enabled":true}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "schedule cron is required",
		},
		{
			name:       "enabled without timezone",
			body:       `{"name":"x","schedule_cron":"0 9 * * *","schedule_enabled":true}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "schedule timezone is required",
		},
		{
			name:       "invalid cron",
			body:       `{"name":"x","schedule_cron":"not-a-cron","schedule_tz":"UTC","schedule_enabled":true}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "invalid schedule cron",
		},
		{
			name:       "invalid timezone",
			body:       `{"name":"x","schedule_cron":"0 9 * * *","schedule_tz":"Moon/Phase1","schedule_enabled":true}`,
			wantStatus: http.StatusBadRequest,
			wantSubstr: "invalid schedule timezone",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/rulesets", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			CreateRuleset(db)(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d want %d, body = %s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tc.wantSubstr) {
				t.Fatalf("body %q should contain %q", rec.Body.String(), tc.wantSubstr)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateRuleset_MethodNotAllowed(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/rulesets", nil)
	rec := httptest.NewRecorder()
	CreateRuleset(db)(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", rec.Code)
	}
}
