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
)

type Post struct{
	ID int
	UserID int
	Title string
	Body string
}

func GetPosts(posts *[]byte){
	err := godotenv.Load()
	HandleError(err)
	resp, err := http.Get(os.Getenv("URLPosts"))
	HandleError(err)
	*posts, err = ioutil.ReadAll(resp.Body)
	HandleError(err)
}

func ParsePost(postStr string, posts *[]Post) {
	var post Post
	var reading map[string]interface{}
	err := json.Unmarshal([]byte(postStr), &reading)
	HandleError(err)
	data := fmt.Sprintf("%v", reading["id"])
	post.ID, err = strconv.Atoi(data)
	HandleError(err)
	data = fmt.Sprintf("%v", reading["userId"])
	post.UserID, err = strconv.Atoi(data)
	HandleError(err)
	post.Title = fmt.Sprintf("%v", reading["title"])
	post.Body = fmt.Sprintf("%v", reading["body"])
	*posts = append(*posts, post)
}

func PostsToDB(posts []Post, db *gorm.DB){
	for i := range posts {
		db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?)",
			posts[i].ID, posts[i].UserID, posts[i].Title, posts[i].Body)
	}
}