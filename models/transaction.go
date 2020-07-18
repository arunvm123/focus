package models

type Transaction interface {
	Begin() DB
	Commit()
	Rollback()
}
