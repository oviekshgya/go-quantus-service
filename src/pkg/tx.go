package pkg

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
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

func UpdateFieldsDynamic(input interface{}) map[string]interface{} {
	v := reflect.ValueOf(input)
	ind := reflect.Indirect(v)
	updateData := make(map[string]interface{})

	for i := 0; i < ind.NumField(); i++ {
		value := ind.Field(i)
		fieldType := ind.Type().Field(i)

		// Ambil tag GORM
		gormTag := fieldType.Tag.Get("gorm")
		if strings.Contains(gormTag, "primaryKey") {
			continue
		}

		columnName := ""
		tags := strings.Split(gormTag, ";")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column:") {
				columnName = strings.TrimPrefix(tag, "column:")
				break
			}
		}

		if columnName == "" {
			continue
		}

		if !value.IsZero() {
			switch v := value.Interface().(type) {
			case *int:
				updateData[columnName] = *v
			case *string:
				updateData[columnName] = *v
			case *float32:
				updateData[columnName] = *v
			case *float64:
				updateData[columnName] = *v
			case *time.Time:
				updateData[columnName] = *v
			default:
				updateData[columnName] = v
			}
		}
	}

	return updateData
}
