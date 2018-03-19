package sqlxpersistence

import (
	"errors"
	"github.com/zibilal/repoman/persistence"
)

type MySqlxRepository struct {
}

func NewMySqlxRepository() *MySqlxRepository {
	return new(MySqlxRepository)
}

func (r *MySqlxRepository) Finder(dbContext persistence.DatabaseContext, output interface{}) error {
	return errors.New("not implemented")
}

func (r *MySqlxRepository) Updater(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (r *MySqlxRepository) Creator(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (r *MySqlxRepository) Deleter(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}
