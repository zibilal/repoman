package persistence

type DatabaseContextBuilder interface {
	Connect(connectionString string) DatabaseContextBuilder
	Build() DatabaseContext
}

type DatabaseContext interface{
	Db() interface{}
	SetTransaction(bool)
	IsTransaction() bool
	Commit() error
	Rollback() error
}
