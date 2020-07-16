package mysql

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arunvm/travail-backend/models"
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
	row1 := models.List{
		Heading:   "Heading1",
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		ID:        "935b639e-2fc9-473d-a2fc-6ecf9562f444",
		TeamID:    uuid.New().String(),
		UserID:    1,
	}

	dbClient, mock, mockdb, err := setup()
	if err != nil {
		t.Errorf("Error initialising mock DB; %v", err)
	}
	defer mockdb.Close()

	db := &Mysql{
		Client: dbClient,
	}

	tables := []struct {
		name   string
		data   *models.List
		input  string
		output *models.List
		err    error
	}{
		{
			name:   "list present",
			data:   &row1,
			input:  row1.ID,
			output: &row1,
			err:    nil,
		},
		{
			name:   "list absent",
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

			list, err := db.getList(tt.input)
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

func TestCreateList(t *testing.T) {
	user := &models.User{
		ID:          1,
		Email:       "random@rhyta.com",
		GoogleOauth: false,
		Name:        "John Doe",
		Password:    "",
		Verified:    true,
	}

	teamID := uuid.New().String()

	tables := []struct {
		name             string
		user             *models.User
		team             *models.Team
		args             models.CreateListArgs
		err              error
		shouldCreateList bool
	}{
		{
			name: "successful list creation",
			args: models.CreateListArgs{
				Heading: "Heading",
				TeamID:  teamID,
			},
			err: nil,
			team: &models.Team{
				ID:             teamID,
				AdminID:        user.ID,
				Archived:       false,
				CreatedAt:      time.Now().Unix(),
				Name:           "Test Team",
				OrganisationID: uuid.New().String(),
			},
			user:             user,
			shouldCreateList: true,
		},
		{
			name: "user not admin",
			user: user,
			team: &models.Team{
				ID:             teamID,
				AdminID:        2, // AdminID value is different from user id
				Archived:       false,
				CreatedAt:      time.Now().Unix(),
				Name:           "Another Test Team",
				OrganisationID: uuid.New().String(),
			},
			args: models.CreateListArgs{
				Heading: "Heading",
				TeamID:  teamID,
			},
			err:              models.UserNotAdminOfTeam,
			shouldCreateList: false,
		},
		{
			name: "team does not exist",
			user: user,
			team: nil,
			args: models.CreateListArgs{
				Heading: "Heading",
				TeamID:  uuid.New().String(),
			},
			err:              gorm.ErrRecordNotFound,
			shouldCreateList: false,
		},
	}

	dbClient, mock, mockdb, err := setup()
	if err != nil {
		t.Errorf("Error initialising mock DB; %v", err)
	}
	defer mockdb.Close()

	db := &Mysql{
		Client: dbClient,
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			if tt.team != nil {
				teamRows := mock.NewRows([]string{"id", "organisation_id", "admin_id", "name", "description", "created_at", "archived"}).
					AddRow(tt.team.ID, tt.team.OrganisationID, tt.team.AdminID, tt.team.Name, tt.team.Description, tt.team.CreatedAt, tt.team.Archived)

				mock.ExpectQuery(".*").WithArgs(tt.args.TeamID).WillReturnRows(teamRows)
			} else {
				mock.ExpectQuery(".*").WillReturnError(gorm.ErrRecordNotFound)
			}

			if tt.shouldCreateList {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists`")).WithArgs(sqlmock.AnyArg(), tt.user.ID, tt.args.TeamID, tt.args.Heading, sqlmock.AnyArg(), false).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			}

			_, err = db.CreateList(&tt.args, tt.user)
			if err != nil {
				if err != tt.err {
					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
				}
			}

			if err == nil {
				assert.Equal(t, tt.err, err)
			}

		})
	}

}

func TestUpdateList(t *testing.T) {
	dbClient, mock, mockdb, err := setup()
	if err != nil {
		t.Errorf("Error initialising mock DB; %v", err)
	}
	defer mockdb.Close()

	db := &Mysql{
		Client: dbClient,
	}

	row1 := models.List{
		Heading:   "Heading1",
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		ID:        "935b639e-2fc9-473d-a2fc-6ecf9562f444",
		TeamID:    uuid.New().String(),
		UserID:    1,
	}

	heading := "heading"
	archived := false
	args := &models.UpdateListArgs{
		Heading:  &heading,
		Archived: &archived,
		ID:       row1.ID,
	}

	tables := []struct {
		name             string
		user             *models.User
		list             *models.List
		args             *models.UpdateListArgs
		err              error
		shouldUpdateList bool
	}{
		{
			name:             "update list successfully",
			list:             &row1,
			args:             args,
			err:              nil,
			shouldUpdateList: true,
			user:             getUser(),
		},
		{
			name:             "list does not exist",
			list:             nil,
			args:             args,
			err:              gorm.ErrRecordNotFound,
			shouldUpdateList: false,
			user:             getUser(),
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			if tt.list != nil {
				rows := mock.NewRows([]string{"id", "user_id", "team_id", "heading", "created_at", "archived"}).
					AddRow(tt.list.ID, tt.list.UserID, tt.list.TeamID, tt.list.Heading, tt.list.CreatedAt, tt.list.Archived)

				mock.ExpectQuery(".*").WillReturnRows(rows)
			} else {
				mock.ExpectQuery(".*").WillReturnError(gorm.ErrRecordNotFound)
			}

			if tt.shouldUpdateList {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
					WithArgs(tt.list.UserID, tt.list.TeamID, tt.args.Heading, tt.list.CreatedAt, tt.args.Archived, tt.args.ID).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			err = db.UpdateList(args, tt.user)
			if err != nil {
				if err != tt.err {
					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
				}
			}

			if err == nil {
				assert.Equal(t, tt.err, err)
			}

		})
	}

}

func getUser() *models.User {
	user := &models.User{
		ID:          1,
		Email:       "johnDoe@rhyta.com",
		GoogleOauth: false,
		Name:        "John Doe",
		Password:    "",
		Verified:    true,
	}

	return user
}
