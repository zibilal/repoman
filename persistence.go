package repoman

type Repository interface{
	Finder
	Updater
	Creator
}

type Finder interface {
	Find(query Query, output interface{}) error
}

type Updater interface {
	Update(input interface{}) (interface{}, error)
}

type Creator interface{
	Create(input interface{}) (interface{}, error)
}

type Query interface {
	Execute(query string, input ...interface{}) (interface{}, error)
}