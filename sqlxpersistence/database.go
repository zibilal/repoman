package sqlxpersistence

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func (d *DbSqlxContextBuilder) Connect(driver, connectionString string) *DbSqlxContextBuilder {
	d.driver = driver
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
	tx            *sqlx.Tx
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
		tx, err := m.db.Beginx()

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
	if m.isTransaction {
		return m.tx
	} else {
		return m.db
	}
}

func (m *DbSqlxContext) Begin() interface{} {
	var err error
	m.tx, err = m.db.Beginx()
	if err != nil {
		panic("failed to begin a transaction, error: " + err.Error())
	}
	m.isTransaction = true

	return m.tx
}

func (m *DbSqlxContext) Commit() error {
	m.isTransaction = false
	return m.tx.Commit()
}

func (m *DbSqlxContext) Rollback() error {
	m.isTransaction = false
	return m.tx.Rollback()
}
