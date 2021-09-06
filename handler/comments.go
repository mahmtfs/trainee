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
	"strings"
)

type Comment struct {
	ID int
	PostID int
	Name string
	Email string
	Body string
}

func CreateChoosePost(c echo.Context) error{
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
	return c.Render(302, "choose-post-create.html", posts)
}

func CreateCommentHandler(c echo.Context) error{
	postID, _ := strconv.Atoi(c.FormValue("id"))
	return c.Render(302, "create-comment.html", postID)
}

func CreateComment (c echo.Context) error{
	var db *gorm.DB
	comment := Comment{}
	GetDB(&db)
	db.Last(&comment)
	comment.ID+=1
	comment.PostID, _ = strconv.Atoi(c.FormValue("postID"))
	comment.Name = c.FormValue("name")
	comment.Email = c.FormValue("email")
	comment.Body = c.FormValue("body")
	db.Exec("INSERT INTO comments VALUES (?, ?, ?, ?, ?);",
		comment.ID, comment.PostID, comment.Name, comment.Email, comment.Body)
	db.Close()
	return c.Redirect(302, "/")
}

func ReadChoosePost(c echo.Context) error{
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
	return c.Render(302, "choose-post-read.html", posts)
}

func ReadCommentHandler(c echo.Context) error{
	postID, _ := strconv.Atoi(c.FormValue("id"))
	type Data struct{
		PostID int
		JSONData string
		XMLData string
	}
	DataHTML := Data{}
	DataHTML.PostID = postID
	return c.Render(302, "read-comment.html", DataHTML)
}

func ReadComment (c echo.Context) error{
	var db *gorm.DB
	comment := Comment{}
	GetDB(&db)
	comment.ID, _ = strconv.Atoi(c.FormValue("id"))
	comment.PostID, _ = strconv.Atoi(c.FormValue("postID"))
	db.Where("id = ? AND post_id= ?", comment.ID, comment.PostID).Find(&comment)
	type Data struct{
		PostID int
		JSONData string
		XMLData string
	}
	DataHTML := Data{}

	JSONData, _ := json.Marshal(comment)
	XMLData, _ := xml.Marshal(comment)
	DataHTML.PostID = comment.PostID
	DataHTML.JSONData = string(JSONData)
	DataHTML.XMLData = string(XMLData)
	db.Close()
	return c.Render(302, "read-comment.html", DataHTML)
}

func FromCommentUpdateToHomeHandler(c echo.Context) error{
	var db *gorm.DB
	GetDB(&db)
	comments := []Comment{}
	comment := Comment{}
	db.First(&comment)
	idFirst := comment.ID
	comment = Comment{}
	db.Last(&comment)
	idLast := comment.ID
	for i:=idFirst; i <= idLast; i++{
		comment = Comment{}
		db.Where("id = ?", i).Find(&comment)
		comments = append(comments, comment)
	}
	comment.ID, _ = strconv.Atoi(c.FormValue("id"))
	comment.PostID, _ = strconv.Atoi(c.FormValue("postID"))
	comment.Body = c.FormValue("body")
	db.Exec("DELETE FROM comments;")
	for i:=0; i < len(comments); i++{
		if comments[i].ID == comment.ID{
			comments[i].Body = comment.Body
		}
		db.Exec("INSERT INTO comments VALUES(?, ?, ?, ?, ?);",
			comments[i].ID, comments[i].PostID, comments[i].Name, comments[i].Email, comments[i].Body)
	}
	db.Close()
	return c.Redirect( 302, "/")
}

func UpdateChoosePost(c echo.Context) error{
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
	return c.Render(302, "choose-post-update.html", posts)
}

func UpdateCommentHandler(c echo.Context) error{
	postID, _ := strconv.Atoi(c.FormValue("id"))
	return c.Render(302, "update-comment-get.html", postID)
}

func UpdateComment (c echo.Context) error{
	var db *gorm.DB
	comment := Comment{}
	GetDB(&db)
	comment.ID, _ = strconv.Atoi(c.FormValue("id"))
	comment.PostID, _ = strconv.Atoi(c.FormValue("postID"))
	db.Where("id = ? AND post_id = ?", comment.ID, comment.PostID).Find(&comment)
	type Data struct {
		ID int
		PostID int
		Body string
	}
	DataHTML := Data{}
	DataHTML.ID = comment.ID
	DataHTML.PostID = comment.PostID
	DataHTML.Body = comment.Body
	db.Close()
	return c.Render(302, "update-comment.html", DataHTML)
}

func DeleteChoosePost(c echo.Context) error{
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
	return c.Render(302, "choose-post-delete.html", posts)
}

func DeleteCommentHandler(c echo.Context) error{
	var db *gorm.DB
	comments := []Comment{}
	comment := Comment{}
	GetDB(&db)
	db.First(&comment)
	idFirst := comment.ID
	comment = Comment{}
	db.Last(&comment)
	idLast := comment.ID
	postID, _ := strconv.Atoi(c.FormValue("id"))
	for i:=idFirst; i <= idLast; i++{
		comment = Comment{}
		db.Where("id = ? AND post_id = ?", i, postID).Find(&comment)
		if comment.ID != 0 {
			comments = append(comments, comment)
		}
	}
	db.Close()
	return c.Render(302, "delete-comment.html", comments)
}

func DeleteComment (c echo.Context) error{
	DeleteID, _ := strconv.Atoi(c.FormValue("option"))
	var db *gorm.DB
	comments := []Comment{}
	comment := Comment{}
	GetDB(&db)
	db.First(&comment)
	idFirst := comment.ID
	comment = Comment{}
	db.Last(&comment)
	idLast := comment.ID
	var deleted = false
	for i:=idFirst; i <= idLast; i++{
		comment = Comment{}
		db.Where("id = ?", i).Find(&comment)
		if deleted == true{
			comment.ID-=1
		}
		if i != DeleteID{
			comments = append(comments, comment)
		} else{
			deleted = true
		}
	}
	db.Exec("DELETE FROM comments;")
	for i:=0; i < len(comments); i++{
		db.Exec("INSERT INTO comments VALUES(?, ?, ?, ?, ?);",
			comments[i].ID, comments[i].PostID, comments[i].Name, comments[i].Email, comments[i].Body)
	}
	db.Close()
	return c.Redirect(302, "/")
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