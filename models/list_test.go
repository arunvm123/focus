package models

// func setup() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
// 	mockdb, mock, err := sqlmock.New()
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	db, err := gorm.Open("mysql", mockdb)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	return db, mock, mockdb, err
// }

// func TestGetList(t *testing.T) {
// 	row1 := List{
// 		Heading:   "Heading1",
// 		Archived:  false,
// 		CreatedAt: time.Now().Unix(),
// 		ID:        "935b639e-2fc9-473d-a2fc-6ecf9562f444",
// 		TeamID:    uuid.New().String(),
// 		UserID:    1,
// 	}

// 	db, mock, mockdb, err := setup()
// 	if err != nil {
// 		t.Errorf("Error initialising mock DB; %v", err)
// 	}
// 	defer mockdb.Close()

// 	tables := []struct {
// 		name   string
// 		data   *List
// 		input  string
// 		output *List
// 		err    error
// 	}{
// 		{
// 			name:   "list present",
// 			data:   &row1,
// 			input:  row1.ID,
// 			output: &row1,
// 			err:    nil,
// 		},
// 		{
// 			name:   "list absent",
// 			data:   nil,
// 			input:  uuid.New().String(),
// 			output: nil,
// 			err:    gorm.ErrRecordNotFound,
// 		},
// 	}

// 	for _, tt := range tables {
// 		t.Run(tt.name, func(t *testing.T) {
// 			query := mock.ExpectQuery(".*")

// 			if tt.data != nil {
// 				rows := mock.NewRows([]string{"id", "user_id", "team_id", "heading", "created_at", "archived"}).
// 					AddRow(tt.data.ID, tt.data.UserID, tt.data.TeamID, tt.data.Heading, tt.data.CreatedAt, tt.data.Archived)

// 				query.WillReturnRows(rows)
// 			} else {
// 				query.WillReturnError(tt.err)
// 			}

// 			list, err := getList(db, tt.input)
// 			if err != nil {
// 				if err != tt.err {
// 					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
// 				}
// 			}

// 			if err == nil {
// 				assert.Equal(t, tt.output, list)
// 			}
// 		})
// 	}

// }

// func TestCreateList(t *testing.T) {
// 	user := &User{
// 		ID:          1,
// 		Email:       "random@rhyta.com",
// 		GoogleOauth: false,
// 		Name:        "John Doe",
// 		Password:    "",
// 		Verified:    true,
// 	}

// 	teamID := uuid.New().String()

// 	tables := []struct {
// 		name             string
// 		user             *User
// 		team             *Team
// 		args             CreateListArgs
// 		err              error
// 		shouldCreateList bool
// 	}{
// 		{
// 			name: "successful list creation",
// 			args: CreateListArgs{
// 				Heading: "Heading",
// 				TeamID:  teamID,
// 			},
// 			err: nil,
// 			team: &Team{
// 				ID:             teamID,
// 				AdminID:        user.ID,
// 				Archived:       false,
// 				CreatedAt:      time.Now().Unix(),
// 				Name:           "Test Team",
// 				OrganisationID: uuid.New().String(),
// 			},
// 			user:             user,
// 			shouldCreateList: true,
// 		},
// 		{
// 			name: "user not admin",
// 			user: user,
// 			team: &Team{
// 				ID:             teamID,
// 				AdminID:        2, // AdminID value is different from user id
// 				Archived:       false,
// 				CreatedAt:      time.Now().Unix(),
// 				Name:           "Another Test Team",
// 				OrganisationID: uuid.New().String(),
// 			},
// 			args: CreateListArgs{
// 				Heading: "Heading",
// 				TeamID:  teamID,
// 			},
// 			err:              userNotAdminOfTeam,
// 			shouldCreateList: false,
// 		},
// 		{
// 			name: "team does not exist",
// 			user: user,
// 			team: nil,
// 			args: CreateListArgs{
// 				Heading: "Heading",
// 				TeamID:  uuid.New().String(),
// 			},
// 			err:              gorm.ErrRecordNotFound,
// 			shouldCreateList: false,
// 		},
// 	}

// 	db, mock, mockdb, err := setup()
// 	if err != nil {
// 		t.Errorf("Error initialising mock DB; %v", err)
// 	}
// 	defer mockdb.Close()

// 	for _, tt := range tables {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.team != nil {
// 				teamRows := mock.NewRows([]string{"id", "organisation_id", "admin_id", "name", "description", "created_at", "archived"}).
// 					AddRow(tt.team.ID, tt.team.OrganisationID, tt.team.AdminID, tt.team.Name, tt.team.Description, tt.team.CreatedAt, tt.team.Archived)

// 				mock.ExpectQuery(".*").WithArgs(tt.args.TeamID).WillReturnRows(teamRows)
// 			} else {
// 				mock.ExpectQuery(".*").WillReturnError(gorm.ErrRecordNotFound)
// 			}

// 			if tt.shouldCreateList {
// 				mock.ExpectBegin()
// 				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists`")).WithArgs(sqlmock.AnyArg(), tt.user.ID, tt.args.TeamID, tt.args.Heading, sqlmock.AnyArg(), false).
// 					WillReturnResult(driver.ResultNoRows)
// 				mock.ExpectCommit()
// 			}

// 			_, err = user.CreateList(db, &tt.args)
// 			if err != nil {
// 				if err != tt.err {
// 					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
// 				}
// 			}

// 			if err == nil {
// 				assert.Equal(t, tt.err, err)
// 			}

// 		})
// 	}

// }

// func TestUpdateList(t *testing.T) {
// 	db, mock, mockdb, err := setup()
// 	if err != nil {
// 		t.Errorf("Error initialising mock DB; %v", err)
// 	}
// 	defer mockdb.Close()

// 	row1 := List{
// 		Heading:   "Heading1",
// 		Archived:  false,
// 		CreatedAt: time.Now().Unix(),
// 		ID:        "935b639e-2fc9-473d-a2fc-6ecf9562f444",
// 		TeamID:    uuid.New().String(),
// 		UserID:    1,
// 	}

// 	heading := "heading"
// 	archived := false
// 	args := &UpdateListArgs{
// 		Heading:  &heading,
// 		Archived: &archived,
// 		ID:       row1.ID,
// 	}

// 	tables := []struct {
// 		name             string
// 		user             *User
// 		list             *List
// 		args             *UpdateListArgs
// 		err              error
// 		shouldUpdateList bool
// 	}{
// 		{
// 			name:             "update list successfully",
// 			list:             &row1,
// 			args:             args,
// 			err:              nil,
// 			shouldUpdateList: true,
// 			user:             getUser(),
// 		},
// 		{
// 			name:             "list does not exist",
// 			list:             nil,
// 			args:             args,
// 			err:              gorm.ErrRecordNotFound,
// 			shouldUpdateList: false,
// 			user:             getUser(),
// 		},
// 	}

// 	for _, tt := range tables {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.list != nil {
// 				rows := mock.NewRows([]string{"id", "user_id", "team_id", "heading", "created_at", "archived"}).
// 					AddRow(tt.list.ID, tt.list.UserID, tt.list.TeamID, tt.list.Heading, tt.list.CreatedAt, tt.list.Archived)

// 				mock.ExpectQuery(".*").WillReturnRows(rows)
// 			} else {
// 				mock.ExpectQuery(".*").WillReturnError(gorm.ErrRecordNotFound)
// 			}

// 			if tt.shouldUpdateList {
// 				mock.ExpectBegin()
// 				mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
// 					WithArgs(tt.list.UserID, tt.list.TeamID, tt.args.Heading, tt.list.CreatedAt, tt.args.Archived, tt.args.ID).WillReturnResult(sqlmock.NewResult(0, 1))
// 				mock.ExpectCommit()
// 			}

// 			err = tt.user.UpdateList(db, args)
// 			if err != nil {
// 				if err != tt.err {
// 					t.Errorf("wrong error behavior %v, wantErr %v", err, tt.err)
// 				}
// 			}

// 			if err == nil {
// 				assert.Equal(t, tt.err, err)
// 			}

// 		})
// 	}

// }

// func getUser() *User {
// 	user := &User{
// 		ID:          1,
// 		Email:       "johnDoe@rhyta.com",
// 		GoogleOauth: false,
// 		Name:        "John Doe",
// 		Password:    "",
// 		Verified:    true,
// 	}

// 	return user
// }
