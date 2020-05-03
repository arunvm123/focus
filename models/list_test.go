package models

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetList(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initialising sql mock: %v", err)
	}
	defer mockdb.Close()

	db, err := gorm.Open("mysql", mockdb)
	if err != nil {
		t.Fatalf("Error initialising gorm mock db: %v", err)
	}

	row1 := List{
		Heading:   "Heading1",
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		ID:        uuid.New().String(),
		TeamID:    uuid.New().String(),
		UserID:    1,
	}

	rows := mock.NewRows([]string{"id", "user_id", "team_id", "heading", "created_at", "archived"}).
		AddRow(row1.ID, row1.UserID, row1.TeamID, row1.Heading, row1.CreatedAt, row1.Archived)

	mock.ExpectQuery("SELECT * FROM `lists` WHERE (archived = false AND id = $2)").WithArgs(row1.ID).WillReturnRows(rows)

	list, err := getList(db, row1.ID)
	if err != nil {
		t.Fatalf("Error when retrieving list: %v", err)
	}

	assert.Equal(t, row1, list)
}
