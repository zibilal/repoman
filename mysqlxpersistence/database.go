package mysqlxpersistence

import "github.com/jmoiron/sqlx"

type DbMySqlxContextBuilder struct {
	connectionString string
	dbContext *DbMySqlxContext
}

func NewDbSqlxContextBuilder() *DbMySqlxContextBuilder {
	db := new(DbMySqlxContextBuilder)

	return db
}

func(d *DbMySqlxContextBuilder) Connect(connectionString string) *DbMySqlxContextBuilder {
	d.connectionString = connectionString

	return d
}

func(d *DbMySqlxContextBuilder) Build() *DbMySqlxContext {
	db, err := sqlx.Connect("mysql", d.connectionString)

	if err != nil {
		panic(err)
	}

	dbContext := NewDbMySqlxContext(db)
	d.dbContext = dbContext

	return dbContext
}

type DbMySqlxContext struct {
	db *sqlx.DB
}

func NewDbMySqlxContext(db *sqlx.DB) *DbMySqlxContext {
	mySql := new(DbMySqlxContext)
	mySql.db = db

	return mySql
}

func(m *DbMySqlxContext) Db() interface{}{
	return m.db
}

func(m *DbMySqlxContext) Tx() interface{}{
	tx, err := m.db.Begin()

	if err != nil {
		panic(err)
	}

	return tx
}


