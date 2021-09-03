package handler

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

func HandleError(err error){
	if err != nil{
		panic(err)
	}
}

func GetDB(db **sql.DB){
	err := godotenv.Load()
	HandleError(err)
	*db, err = sql.Open(os.Getenv("DbDriver"), os.Getenv("DataSource"))
	HandleError(err)
	dbr := *db
	err = dbr.Ping()
	HandleError(err)
	_, err = dbr.Query("DELETE FROM posts;")
	HandleError(err)
	_, err = dbr.Query("DELETE FROM comments;")
	HandleError(err)
}
