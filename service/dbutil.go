package service

import (
	"github.com/jmoiron/sqlx"
)

type DBUtil struct {
	db *sqlx.DB
}

func NewDBUtil(dbInstance *sqlx.DB) *DBUtil {
	return &DBUtil{
		db: dbInstance,
	}
}

func (dbUtil *DBUtil) RunTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := dbUtil.db.Beginx()
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
