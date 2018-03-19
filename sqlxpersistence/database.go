package sqlxpersistence

import "github.com/jmoiron/sqlx"

type DbSqlxContextBuilder struct {
	connectionString string
	driver           string
	dbContext        *DbSqlxContext
}

func NewDbSqlxContextBuilder() *DbSqlxContextBuilder {
	db := new(DbSqlxContextBuilder)

	return db
}

func (d *DbSqlxContextBuilder) Connect(connectionString string) *DbSqlxContextBuilder {
	d.connectionString = connectionString

	return d
}

func (d *DbSqlxContextBuilder) Build() *DbSqlxContext {
	db, err := sqlx.Connect(d.driver, d.connectionString)

	if err != nil {
		panic(err)
	}

	dbContext := NewDbMySqlxContext(db)
	d.dbContext = dbContext

	return dbContext
}

type DbSqlxContext struct {
	db *sqlx.DB
}

func NewDbMySqlxContext(db *sqlx.DB) *DbSqlxContext {
	theDb := new(DbSqlxContext)
	theDb.db = db

	return theDb
}

func (m *DbSqlxContext) Db() interface{} {
	return m.db
}
