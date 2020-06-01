package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/assert.v1"
)

func setup() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := gorm.Open("mysql", mockdb)
	if err != nil {
		return nil, nil, nil, err
	}

	return db, mock, mockdb, err
}

func TestGetList(t *testing.T) {
	row1 := List{
		Heading:   "Heading1",
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		ID:        "935b639e-2fc9-473d-a2fc-6ecf9562f444",
		TeamID:    uuid.New().String(),
		UserID:    1,
	}

	db, mock, mockdb, err := setup()
	if err != nil {
		t.Errorf("Error initialising mock DB; %v", err)
	}
	defer mockdb.Close()

	tables := []struct {
		name   string
		data   *List
		input  string
		output *List
		err    error
	}{
		{
			name:   "List present",
			data:   &row1,
			input:  row1.ID,
			output: &row1,
			err:    nil,
		},
		{
			name:   "List absent",
			data:   nil,
			input:  uuid.New().String(),
			output: nil,
			err:    gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			query := mock.ExpectQuery(".*")

			if tt.data != nil {
				rows := mock.NewRows([]string{"id", "user_id", "team_id", "heading", "created_at", "archived"}).
					AddRow(tt.data.ID, tt.data.UserID, tt.data.TeamID, tt.data.Heading, tt.data.CreatedAt, tt.data.Archived)

				query.WillReturnRows(rows)
			} else {
				query.WillReturnError(tt.err)
			}

			list, err := getList(db, tt.input)
			if err != nil {
				if err != tt.err {
					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
				}
			}

			if err == nil {
				assert.Equal(t, tt.output, list)
			}
		})
	}

}
