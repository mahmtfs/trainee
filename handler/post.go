package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

func PostCreateHandler(c echo.Context) error{
	return c.Render(http.StatusOK, "create-post.html", nil)
}

func CreatePost(c echo.Context) error{
	fmt.Println("Starting...")
	post := Post{}
	var db *gorm.DB
	GetDB(&db)
	db.Last(&post)
	post.Title = c.FormValue("title")
	post.Body = c.FormValue("body")
	post.ID += 1
	post.UserID = 7
	db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?)",
	post.ID, post.UserID, post.Title, post.Body)
	db.Close()
	return c.Redirect(http.StatusFound, "/")
}

func PostReadHandler(c echo.Context) error{
	return c.Render(302, "read-post.html", nil)
}

func ReadPost(c echo.Context) error{
	var db *gorm.DB
	post := Post{}
	GetDB(&db)
	post.ID, _ = strconv.Atoi(c.FormValue("id"))
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
	return c.Render(302, "read-post.html", DataHTML)
}

func FromPostUpdateToHomeHandler(c echo.Context) error{
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
	post.ID, _ = strconv.Atoi(c.FormValue("id"))
	post.Title = c.FormValue("title")
	post.Body = c.FormValue("body")
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
	return c.Redirect(302, "/")
}

func PostUpdateHandler(c echo.Context) error{
	return c.Render(302, "update-post-get.html", nil)
}

func UpdatePost(c echo.Context) error{
	var db *gorm.DB
	post := Post{}
	GetDB(&db)
	post.ID, _ = strconv.Atoi(c.FormValue("id"))
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
	return c.Render(302, "update-post.html", DataHTML)
}

func DeletePostHandler(c echo.Context) error{
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
	return c.Render(302, "delete-post.html", posts)
}

func DeletePost (c echo.Context) error{
	DeteletID, _ := strconv.Atoi(c.FormValue("option"))
	var db *gorm.DB
	posts := []Post{}
	post := Post{}
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
	return c.Redirect(302, "/")
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