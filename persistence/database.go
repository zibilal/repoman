package persistence

type DatabaseContextBuilder interface {
	Connect(connectionString string) DatabaseContextBuilder
	Build() DatabaseContext
}

type DatabaseContext interface{
	Db() interface{}
	Tx() interface{}
}
