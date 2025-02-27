package util

import (
	"github.com/jmoiron/sqlx"
)

type TxHelper struct {
	db *sqlx.DB
}

func NewTxHelper(dbInstance *sqlx.DB) *TxHelper {
	return &TxHelper{
		db: dbInstance,
	}
}

func (txHelper *TxHelper) RunTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := txHelper.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = fn(tx)
	return err
}
