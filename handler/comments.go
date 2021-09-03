package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Comment struct {
	ID int
	PostID int
	Name string
	Email string
	Body string
}

func CommentsProcessing(id int, synch chan int, db *gorm.DB){
	var commentStrings []string
	var comments []Comment
	ch := make(chan int)
	err := godotenv.Load()
	HandleError(err)
	resp, err := http.Get(os.Getenv("URLComments") + strconv.Itoa(id))
	HandleError(err)
	commentsBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err)
	commentStrings = strings.Split(string(commentsBytes), "}")
	for i := 0; i < len(commentStrings) - 1; i++{
		ParseComment(commentStrings[i][1:] + "}", &comments)
	}
	go func(){
		CommentsToDB(comments, db)
		ch <- 1
	}()
	<- ch
	synch <- 1
}

func ParseComment(commentStr string, comments *[]Comment){
	var comment Comment
	var reading map[string]interface{}
	err := json.Unmarshal([]byte(commentStr), &reading)
	HandleError(err)
	data := fmt.Sprintf("%v", reading["id"])
	comment.ID, err = strconv.Atoi(data)
	HandleError(err)
	data = fmt.Sprintf("%v", reading["postId"])
	comment.PostID, err = strconv.Atoi(data)
	HandleError(err)
	comment.Name = fmt.Sprintf("%v", reading["name"])
	comment.Body = fmt.Sprintf("%v", reading["body"])
	comment.Email = fmt.Sprintf("%v", reading["email"])
	*comments = append(*comments, comment)
}

func CommentsToDB(comments []Comment, db *gorm.DB){
	for i := range comments {
		db.Exec("INSERT INTO comments VALUES (?, ?, ?, ?, ?)",
			comments[i].ID, comments[i].PostID, comments[i].Name, comments[i].Email, comments[i].Body)
	}
}