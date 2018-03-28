package repoman

import (
	"errors"
	"fmt"
	"github.com/zibilal/repoman/persistence"
)

type Repository interface {
	Finder
	Updater
	Creator
	QueryMapper
}

type QueryMapper interface {
	AddQuery(string, string)
	GetQuery(string) string
}

type Finder interface {
	Find(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{},  error)
}

type Updater interface {
	Update(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error)
}

type Creator interface {
	Create(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error)
}

type Deleter interface {
	Delete(dbContext persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error)
}

type RepoFunc func(context persistence.DatabaseContext, queryName string, data ...interface{}) (interface{}, error)

type ExecuteQueryOutput struct {
	Output map[int]interface{}
}

type ExecuteQueryInput struct {
	Data      []interface{}
	QueryName string
	Handler   RepoFunc
}

func ExecuteQueryHandlerFunc(dbContext persistence.DatabaseContext, inputs ...ExecuteQueryInput) (*ExecuteQueryOutput, error) {

	var (
		tmp       interface{}
		err       error
		tmpError  error
		message   string
		outputMap map[int]interface{}
	)

	outputMap = make(map[int]interface{})

	for idx, input := range inputs {
		tmp, err = input.Handler(dbContext, input.QueryName, input.Data...)

		if dbContext.IsTransaction() && err != nil {
			tmpError = dbContext.Rollback()

			message = "unable to finish process due to %s"

			if tmpError != nil {
				message += message + ", failed try to rollback due to %s"
			}

			return nil, errors.New(message)
		}

		outputMap[idx] = tmp
	}

	err = dbContext.Commit()

	if err != nil {
		return nil, fmt.Errorf("process finished successfully, but commit is failed due to %s", err.Error())
	}

	return &ExecuteQueryOutput{
		outputMap,
	}, nil
}
