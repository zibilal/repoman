package repoman

import "github.com/zibilal/repoman/persistence"

type Repository interface{
	Finder
	Updater
	Creator
}

type Finder interface {
	Find(dbContext persistence.DatabaseContext, output interface{}) error
}

type Updater interface {
	Update(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error)
}

type Creator interface{
	Create(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error)
}

type Deleter interface {
	Delete(dbContext persistence.DatabaseContext, input interface{}) (interface{}, error)
}

type Query interface {
	Execute(query string, input ...interface{}) (interface{}, error)
}