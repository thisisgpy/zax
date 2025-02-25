package service

import (
	"github.com/jmoiron/sqlx"
)

type TxWrapper struct {
	db *sqlx.DB
}

func NewTxWrapper(dbInstance *sqlx.DB) *TxWrapper {
	return &TxWrapper{
		db: dbInstance,
	}
}

func (txWrapper *TxWrapper) RunTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := txWrapper.db.Beginx()
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
