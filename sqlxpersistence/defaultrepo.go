package sqlxpersistence

import (
	"github.com/zibilal/repoman/persistence"
	"errors"
	"github.com/zibilal/repoman/query"
	"github.com/jmoiron/sqlx"
	"fmt"
)

type DefaultRepo struct {}

func NewDefaultRepo() *DefaultRepo {
	return new(DefaultRepo)
}

func (d *DefaultRepo) Find(dbContext persistence.DatabaseContext, data ...interface{}) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid state, please provide at least one struct query and data destination")
	}

	sQuery, err := query.ComposeQuery(query.QuerySelect, data[0], true)

	if err != nil {
		return nil, err
	}

	if dbContext.IsTransaction() {
		tx := dbContext.Db().(*sqlx.Tx)
		if len(data) > 2 {
			err := tx.Select(data[1], sQuery, data[2:])
			if err != nil {
				return nil, err
			}
		} else {
			err := tx.Select(data[1], sQuery)
			fmt.Println("The error", err)
			if err != nil {
				return nil, err
			}
		}
	} else {
		db := dbContext.Db().(*sqlx.DB)
		if len(data) > 2 {
			err := db.Select(data[1], sQuery, data[2:])
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Select(data[1], sQuery)
			if err != nil {
				return nil, err
			}
		}
	}

	return data[1], nil
}

func (d *DefaultRepo) Update(dbContext persistence.DatabaseContext, data ...interface{}) (interface{}, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid state, please provide at least one struct data")
	}

	sQuery, err := query.ComposeQuery(query.QueryUpdate, data[0], true)

	if err != nil {
		return nil, err
	}

	if dbContext.IsTransaction() {
		tx := dbContext.Db().(*sqlx.Tx)
		if len(data) > 1 {
			result, err := tx.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := tx.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	} else {
		db := dbContext.Db().(*sqlx.DB)
		if len(data) > 1 {
			result, err := db.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := db.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}

func (d *DefaultRepo) Insert(dbContext persistence.DatabaseContext, data ...interface{}) (interface{}, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid state, please provide at least one struct data")
	}

	sQuery, err := query.ComposeQuery(query.QueryInsert, data[0], true)
	if err != nil {
		return nil, err
	}

	if dbContext.IsTransaction() {
		tx := dbContext.Db().(*sqlx.Tx)
		if len(data) > 1 {
			result, err := tx.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := tx.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, err
		}
	} else {
		db := dbContext.Db().(*sqlx.Tx)
		if len(data) > 1 {
			result, err := db.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := db.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}

func (d *DefaultRepo) Delete(dbContext persistence.DatabaseContext, data ...interface{}) (interface{}, error) {
	if len(data) < 1 {
		return nil, errors.New("invalid state, please provide at least on struct data")
	}

	sQuery, err := query.ComposeQuery(query.QueryDelete, data[0], true)
	if err != nil {
		return nil, err
	}

	if dbContext.IsTransaction() {
		tx := dbContext.Db().(*sqlx.Tx)
		if len(data) > 1 {
			result, err := tx.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := tx.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	} else {
		db := dbContext.Db().(*sqlx.DB)
		if len(data) > 1 {
			result, err := db.Exec(sQuery, data[1:])
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := db.Exec(sQuery)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
}
