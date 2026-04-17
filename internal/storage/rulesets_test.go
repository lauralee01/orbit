package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateRuleset_EmptyCronDisablesSchedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO rulesets`).
		WithArgs("n", "", "", "UTC", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	id, err := CreateRuleset(context.Background(), db, "n", "", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if id != 10 {
		t.Fatalf("id = %d want 10", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
