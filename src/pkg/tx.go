package pkg

import (
	"fmt"
	"gorm.io/gorm"
)

func WithTransaction(db *gorm.DB, fn func(tz *gorm.DB) (interface{}, error)) (interface{}, error) {

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			_ = tx.Rollback()
		} else {
			cerr := tx.Commit().Error
			if cerr != nil {
				tx.Error = fmt.Errorf("error committing transaction: %v", cerr)
			}
		}
	}()
	res, err := fn(tx)
	if err != nil {
		tx.Error = err
		return nil, err
	}

	return res, nil
}

type TxResult struct {
	Data     interface{}
	Err      error
	Rollback bool
}

func WithTransactionV2(db *gorm.DB, fn func(tz *gorm.DB) *TxResult) (interface{}, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			_ = tx.Rollback()
		} else {
			cerr := tx.Commit().Error
			if cerr != nil {
				tx.Error = fmt.Errorf("error committing transaction: %v", cerr)
			}
		}
	}()

	res := fn(tx)
	if res.Rollback {
		tx.Error = res.Err
		return nil, res.Err
	}

	return res.Data, res.Err
}
