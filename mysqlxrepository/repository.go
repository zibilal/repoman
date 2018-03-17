package mysqlxrepository

import (
	"github.com/zibilal/repoman/mysqlxpersistence"
	"errors"
)

type MySqlxRepository struct {
	dbContext mysqlxpersistence.DbMySqlxContextBuilder
}

func NewMySqlxRepository(dbContext mysqlxpersistence.DbMySqlxContext) (*MySqlxRepository, error) {
	return nil, errors.New("unimplemented")
}


