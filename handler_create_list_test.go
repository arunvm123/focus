package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arunvm/focus/models"
	"github.com/arunvm/focus/models/mockdb"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateList(t *testing.T) {
	tables := []struct {
		name               string
		args               *models.CreateListArgs
		user               *models.User
		expectedStatusCode int
		responseErr        error
		expectedResponse   *models.List
		dbErr              error
	}{
		{
			name: "Successful List Creation",
			args: &models.CreateListArgs{
				Heading: "Heading",
				TeamID:  "randomteamID",
			},
			user:               &models.User{},
			expectedStatusCode: http.StatusOK,
			responseErr:        nil,
			expectedResponse:   &models.List{},
			dbErr:              nil,
		},
		{
			name:               "User missing in context",
			args:               nil,
			user:               nil,
			expectedStatusCode: http.StatusBadRequest,
			responseErr:        errors.New("Error fetching user"),
			expectedResponse:   nil,
			dbErr:              nil,
		},
		{
			name:               "Missing request arguments",
			args:               nil,
			user:               &models.User{},
			expectedStatusCode: http.StatusBadRequest,
			responseErr:        errors.New("Request body not properly formatted"),
			expectedResponse:   nil,
			dbErr:              nil,
		},
		{
			name: "CreateList returns DB error",
			args: &models.CreateListArgs{
				Heading: "Heading",
				TeamID:  "randomteamID",
			},
			user:               &models.User{},
			expectedStatusCode: http.StatusInternalServerError,
			responseErr:        errors.New("Error when creating list"),
			expectedResponse:   &models.List{},
			dbErr:              errors.New("Error when creating list"),
		},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			s := newServer()
			mockDB := mockdb.New()
			s.db = mockDB

			if tt.expectedResponse != nil {
				mockDB.On("CreateList", tt.args, tt.user).Return(tt.expectedResponse, tt.dbErr)
			}

			b, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("Error marshalling data: %v", err)
			}

			rr := httptest.NewRecorder()

			ctx, r := gin.CreateTestContext(rr)

			if tt.user != nil {
				r.Use(func(c *gin.Context) {
					c.Keys = make(map[string]interface{})
					c.Keys["user"] = tt.user
				})
			}

			ctx.Request, err = http.NewRequest("POST", "/create/list", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}

			r.POST("/create/list", s.createList)
			r.ServeHTTP(rr, ctx.Request)

			mockDB.AssertExpectations(t)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			if tt.responseErr != nil {
				assert.Equal(t, tt.responseErr.Error(), strings.Trim(strings.ReplaceAll(rr.Body.String(), "\"", ""), "\t \n"))
			}
		})
	}
}
