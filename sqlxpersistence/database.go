package sqlxpersistence

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"database/sql"
)

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
	db            *sqlx.DB
	tx            *sql.Tx
	isTransaction bool
}

func NewDbMySqlxContext(db *sqlx.DB) *DbSqlxContext {
	theDb := new(DbSqlxContext)
	theDb.db = db

	return theDb
}

func (m *DbSqlxContext) SetTransaction(isTransaction bool) {
	m.isTransaction = isTransaction

	if m.isTransaction {
		tx, err := m.db.Begin()

		if err != nil {
			tmp := fmt.Sprintf("failed to beginning transaction, due to %s", err.Error())
			panic(tmp)
		}

		m.tx = tx
	}
}

func (m *DbSqlxContext) IsTransaction() bool {
	return m.isTransaction
}

func (m *DbSqlxContext) Db() interface{} {
	return m.db
}

func (m *DbSqlxContext) Commit() error {
	return m.tx.Commit()
}

func (m *DbSqlxContext) Rollback() error {
	return m.tx.Rollback()
}
