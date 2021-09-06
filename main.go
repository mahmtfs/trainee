package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"strings"
	"traineeProject/handler"
)

type Renderer struct {
	template *template.Template
	debug bool
	location string
}

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
	e := echo.New()
	e.Renderer = NewRenderer("templates/*.html", true)
	e.GET("/", handler.HomeHandler)
	e.GET("/post/create/", handler.PostCreateHandler)
	e.GET("/post/created/", handler.CreatePost)
	e.GET("/post/read/", handler.PostReadHandler)
	e.GET("/post/reaD/", handler.ReadPost)
	e.GET("/post/update/", handler.PostUpdateHandler)
	e.GET("/post/updated/", handler.UpdatePost)
	e.GET("/post/updateD/", handler.FromPostUpdateToHomeHandler)
	e.GET("/post/delete/", handler.DeletePostHandler)
	e.GET("/post/deleted/", handler.DeletePost)
	e.GET("/comment/create/choosepost/", handler.CreateChoosePost)
	e.GET("/comment/create/", handler.CreateCommentHandler)
	e.GET("/comment/created/", handler.CreateComment)
	e.GET("/comment/read/choosepost/", handler.ReadChoosePost)
	e.GET("/comment/read/", handler.ReadCommentHandler)
	e.GET("/comment/reaD/", handler.ReadComment)
	e.GET("/comment/update/choosepost/", handler.UpdateChoosePost)
	e.GET("/comment/update/", handler.UpdateCommentHandler)
	e.GET("comment/updated/", handler.UpdateComment)
	e.GET("comment/updateD/", handler.FromCommentUpdateToHomeHandler)
	e.GET("/comment/delete/choosepost/", handler.DeleteChoosePost)
	e.GET("/comment/delete/", handler.DeleteCommentHandler)
	e.GET("/comment/deleted/", handler.DeleteComment)
	e.Start(":8000")
}

func (t *Renderer) ReloadTemplates(){
	t.template = template.Must(template.ParseGlob(t.location))
}

func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	if t.debug{
		t.ReloadTemplates()
	}
	return t.template.ExecuteTemplate(w, name, data)
}

func NewRenderer(location string, debug bool) *Renderer{
	tpl := new(Renderer)
	tpl.location = location
	tpl.debug = debug

	tpl.ReloadTemplates()
	return tpl
}