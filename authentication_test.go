package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSignup(t *testing.T) {
	server := server{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	server.db, err = gorm.Open("mysql", db)
	if err != nil {
		t.Error("Error initialising gorm")
	}

	mock.ExpectExec("SELECT count(*) from users").WithArgs("test@gmail.com").WillReturnResult(sqlmock.NewResult(1, 0))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, err := json.Marshal(loginRequest{
		Email:    "test@gmail.com",
		Password: "password",
	})
	if err != nil {
		t.Error("Error marshalling request body")
	}

	c.Request = new(http.Request)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

	server.signup(c)
}
