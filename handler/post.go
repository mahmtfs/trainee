package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"html/template"
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

func PostCreateHandler(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/create-post.html"))
	t.Execute(w,nil)
}

func (post Post) Create(w http.ResponseWriter, r *http.Request){
	var db *gorm.DB
	GetDB(&db)
	db.Last(&post)
	err := r.ParseForm()
	HandleError(err)
	post.Title = r.Form.Get("title")
	post.Body = r.Form.Get("body")
	post.ID+=1
	post.UserID = 7
	db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?)",
		post.ID, post.UserID, post.Title, post.Body)
	db.Close()
	http.Redirect(w, r, "/", 302)
}

func PostReadHandler(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/read-post.html"))
	t.Execute(w,nil)
}

func (post Post) Read(w http.ResponseWriter, r *http.Request){
	var db *gorm.DB
	t := template.Must(template.ParseFiles("templates/read-post.html"))
	GetDB(&db)
	err := r.ParseForm()
	HandleError(err)
	post.ID, _ = strconv.Atoi(r.Form.Get("id"))
	db.Where("id = ?", post.ID).Find(&post)
	type Data struct{
		JSONData string
		XMLData string
	}
	DataHTML := Data{}
	JSONData, _ := json.Marshal(post)
	XMLData, _ := xml.Marshal(post)
	DataHTML.JSONData = string(JSONData)
	DataHTML.XMLData = string(XMLData)
	db.Close()
	t.Execute(w, DataHTML)
}

func FromPostUpdateToHomeHandler(w http.ResponseWriter, r *http.Request){
	var db *gorm.DB
	GetDB(&db)
	posts := []Post{}
	post := Post{}
	db.First(&post)
	idFirst := post.ID
	post = Post{}
	db.Last(&post)
	idLast := post.ID
	for i:=idFirst; i <= idLast; i++{
		post = Post{}
		db.Where("id = ?", i).Find(&post)
		posts = append(posts, post)
	}
	r.ParseForm()
	post.ID, _ = strconv.Atoi(r.Form.Get("id"))
	post.Title = r.Form.Get("title")
	post.Body = r.Form.Get("body")
	db.Exec("DELETE FROM posts;")
	for i:=0; i < len(posts); i++{
		if posts[i].ID == post.ID{
			posts[i].Title = post.Title
			posts[i].Body = post.Body
		}
		db.Exec("INSERT INTO posts VALUES(?, ?, ?, ?);",
			posts[i].ID, posts[i].UserID, posts[i].Title, posts[i].Body)
	}
	db.Close()
	http.Redirect(w, r, "/", 302)
}

func PostUpdateHandler(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/update-post-get.html"))
	t.Execute(w, nil)
}

func (post Post) Update(w http.ResponseWriter, r *http.Request){
	var db *gorm.DB
	GetDB(&db)
	t := template.Must(template.ParseFiles("templates/update-post.html"))
	r.ParseForm()
	post.ID, _ = strconv.Atoi(r.Form.Get("id"))
	db.Where("id = ?", post.ID).Find(&post)
	type Data struct {
		ID int
		Title string
		Body string
	}
	DataHTML := Data{}
	DataHTML.ID = post.ID
	DataHTML.Title = post.Title
	DataHTML.Body = post.Body
	db.Close()
	t.Execute(w, DataHTML)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/delete-post.html"))
	var db *gorm.DB
	posts := []Post{}
	post := Post{}
	GetDB(&db)
	db.First(&post)
	idFirst := post.ID
	post = Post{}
	db.Last(&post)
	idLast := post.ID
	for i:=idFirst; i <= idLast; i++{
		post = Post{}
		db.Where("id = ?", i).Find(&post)
		posts = append(posts, post)
	}
	db.Close()
	t.Execute(w, posts)
}

func (post Post) Delete(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	DeteletID, _ := strconv.Atoi(r.Form.Get("option"))
	var db *gorm.DB
	posts := []Post{}
	post = Post{}
	GetDB(&db)
	db.First(&post)
	idFirst := post.ID
	post = Post{}
	db.Last(&post)
	idLast := post.ID
	var deleted = false
	for i:=idFirst; i <= idLast; i++{
		post = Post{}
		db.Where("id = ?", i).Find(&post)
		if deleted == true{
			post.ID-=1
		}
		if i != DeteletID{
			posts = append(posts, post)
		} else{
			deleted = true
		}
	}
	db.Exec("DELETE FROM posts;")
	for i:=0; i < len(posts); i++{
		db.Exec("INSERT INTO posts VALUES(?, ?, ?, ?);",
			posts[i].ID, posts[i].UserID, posts[i].Title, posts[i].Body)
	}
	db.Close()
	http.Redirect(w, r, "/", 302)
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