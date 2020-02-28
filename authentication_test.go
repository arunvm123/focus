package main

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSignup(t *testing.T) {
	server := server{}

	db, mock, err := sqlmock.New()

	server.db, err = gorm.Open("mysql", db)
	if err != nil {
		t.Error("Error")
	}
	mock.ExpectBegin()
}
