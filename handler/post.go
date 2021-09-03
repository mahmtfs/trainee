package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func PostsToDB(posts []Post, db *sql.DB){
	for i := range posts {
		_, err := db.Query("INSERT INTO posts VALUES (?, ?, ?, ?)",
			posts[i].ID, posts[i].UserID, posts[i].Title, posts[i].Body)
		HandleError(err)
	}
}