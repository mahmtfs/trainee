package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"traineeProject/handler"
)

func main(){
	start := time.Now()
	var db *sql.DB
	var postsBytes []byte
	var postsStrings []string
	var posts []handler.Post
	ch := make(chan int)
	handler.GetDB(&db)
	handler.GetPosts(&postsBytes)
	postsStrings = strings.Split(string(postsBytes), "}")
	for i := 0; i < len(postsStrings) - 1; i++{
		handler.ParsePost(postsStrings[i][1:] + "}", &posts)
	}
	handler.PostsToDB(posts, db)
	for i := range posts {
		go handler.CommentsProcessing(posts[i].ID, ch, db)
		<-ch
	}
	fmt.Println(time.Since(start))
}
