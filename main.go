package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"traineeProject/handler"
)

func main(){
	var db *gorm.DB
	var postsBytes []byte
	var postsStrings []string
	var posts []handler.Post
	ch := make(chan int)
	handler.GetDB(&db)
	db.Exec("DELETE FROM posts;")
	db.Exec("DELETE FROM comments;")
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
	db.Close()
	fmt.Println("Starting the server...")
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/post/create/", handler.PostCreateHandler)
	http.HandleFunc("/post/created/", handler.Post{}.Create)
	http.HandleFunc("/post/read/", handler.PostReadHandler)
	http.HandleFunc("/post/reaD/", handler.Post{}.Read)
	http.HandleFunc("/post/update/", handler.PostUpdateHandler)
	http.HandleFunc("/post/updated/", handler.Post{}.Update)
	http.HandleFunc("/post/updateD/", handler.FromPostUpdateToHomeHandler)
	http.HandleFunc("/post/delete/", handler.DeletePostHandler)
	http.HandleFunc("/post/deleted/", handler.Post{}.Delete)
	http.HandleFunc("/comment/create/choosepost/", handler.CreateChoosePost)
	http.HandleFunc("/comment/create/", handler.CreateCommentHandler)
	http.HandleFunc("/comment/created/", handler.Comment{}.Create)
	http.HandleFunc("/comment/read/choosepost/", handler.ReadChoosePost)
	http.HandleFunc("/comment/read/", handler.ReadCommentHandler)
	http.HandleFunc("/comment/reaD/", handler.Comment{}.Read)
	http.HandleFunc("/comment/update/choosepost/", handler.UpdateChoosePost)
	http.HandleFunc("/comment/update/", handler.UpdateCommentHandler)
	http.HandleFunc("/comment/updated/", handler.Comment{}.Update)
	http.HandleFunc("/comment/updateD/", handler.FromCommentUpdateToHomeHandler)
	http.HandleFunc("/comment/delete/choosepost/", handler.DeleteChoosePost)
	http.HandleFunc("/comment/delete/", handler.DeleteCommentHandler)
	http.HandleFunc("/comment/deleted/", handler.Comment{}.Delete)
	http.ListenAndServe(":8000", nil)
}
