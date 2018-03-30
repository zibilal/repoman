package main

import (
	"database/sql"
	"fmt"
	"github.com/zibilal/repoman/sqlxpersistence"
	"log"
	"time"
)

type User struct {
	Id          int    `query:"id,col,primary#default,where" db:"id"`
	Username    string `query:"username,col" db:"username"`
	Firstname   string `query:"firstname,col" db:"firstname"`
	Lastname    string `query:"lastname,col" db:"lastname"`
	Email       string `query:"email,col" db:"email"`
	CreatedTime string `query:"created_time,col" db:"created_time"`
	UpdatedTime string `query:"updated_time,col" db:"updated_time"`
	CreatedBy   int    `query:"created_by,col" db:"created_by"`
	UpdatedBy   int    `query:"updated_by,col" db:"updated_by"`
}

func main() {
	builder := sqlxpersistence.NewDbSqlxContextBuilder()
	dbContext := builder.Connect("mysql", "root:mytestdb@(localhost:3306)/mytestdb").Build()

	fmt.Println("The now", time.Now().Format("2006-01-02 15:04:05"))

	repo := sqlxpersistence.NewDefaultRepo()

	dbContext.Begin()

	users := []User {
		{
			Id:          112235,
			Username:    "admintest",
			Firstname:   "Admin",
			Lastname:    "Test",
			Email:       "admin.test@example.com",
			CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedTime: time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:   123,
			UpdatedBy:   123,
		},
		{
			Id:          112236,
			Username:    "user1test",
			Firstname:   "User",
			Lastname:    "One",
			Email:       "user1.test@example.com",
			CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedTime: time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:   124,
			UpdatedBy:   124,
		},
	}

	for _, user := range users {
		iresult, err := repo.Insert(dbContext, user)

		if err != nil {
			err2 := dbContext.Rollback()
			if err2 != nil {
				log.Fatal("Failed trying to rollback:", err2)
			}
			log.Fatal("Failed due to, error: ", err)

			return
		}

		result := iresult.(sql.Result)
		insertedId, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()
		log.Println("Insert: ", insertedId, rowsAffected)
	}

	fmt.Println("The now", time.Now().Format("2006-01-02 15:04:05"))

	selectResult := []User{}

	user := User{
		Id: 112233,
	}

	_, err := repo.Find(dbContext, user, &selectResult)
	if err != nil {
		log.Fatal("Failed try to get the user data, due to ", err)

		return
	}

	fmt.Println("Select data")
	for _, user := range selectResult {
		fmt.Printf("Id:%d\tUsername:%s\tFirstname:%s\tLastname:%s\tEmail:%s\tCreatedTime:%s\tUpdatedTime:%s\tCreatedBy:%d\tUpdatedBy:%d\n",
			user.Id, user.Username, user.Firstname, user.Lastname, user.Email, user.CreatedTime, user.UpdatedTime, user.CreatedBy, user.UpdatedBy)
	}

	iresult, err := repo.Update(dbContext, selectResult[0])
	if err != nil {
		log.Fatal("Failed try to update user data")
	}

	fmt.Println("Updated", iresult)

	iresult, err = repo.Delete(dbContext, user)
	if err != nil {
		log.Fatal("Failed try to delete user id 112233")
	}

	fmt.Println("Deleted id 112233", iresult)

	err = dbContext.Commit()
	if err != nil {
		log.Fatal("Failed try to commit the transaction")
	}

	log.Println("Insert transaction is committed")
}
