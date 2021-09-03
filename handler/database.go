package handler

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
)

func HandleError(err error){
	if err != nil{
		panic(err)
	}
}

func GetDB(db **gorm.DB){
	err := godotenv.Load()
	HandleError(err)
	*db, err = gorm.Open(os.Getenv("DbDriver"), os.Getenv("DataSource"))
	HandleError(err)
	dbr := *db
	dbr.Exec("DELETE FROM posts;")
	dbr.Exec("DELETE FROM comments;")
}
