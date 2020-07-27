package mysql

import "github.com/arunvm/travail-backend/models"

func (db *Mysql) Begin() models.Transaction {
	return &Mysql{
		Client: db.Client.Begin(),
	}
}

func (tx *Mysql) Commit() {
	tx.Client.Commit()
}

func (tx *Mysql) Rollback() {
	tx.Client.Rollback()
}
