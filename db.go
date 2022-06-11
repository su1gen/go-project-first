package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Post struct {
	Id     uint   `json:"id"`
	UserId uint   `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var postsList []Post
var wgGetComments sync.WaitGroup
var wgSaveComments sync.WaitGroup

func getUserPosts(userId int) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%v", userId)
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &postsList)
}

func savePostToDB(post Post, db *sql.DB) bool {
	stmt, err := db.Prepare("INSERT INTO posts(ID, user_id, title, body) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return false
	}
	res, err := stmt.Exec(post.Id, post.UserId, post.Title, post.Body)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(res)
	return true
}

func getPostComments(postId uint, db *sql.DB) {
	var commentsList []Comment
	defer wgGetComments.Done()
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%v", postId)
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &commentsList)

	for _, comment := range commentsList {
		wgSaveComments.Add(1)
		go saveCommentToDB(comment, db)
	}

}

func saveCommentToDB(comment Comment, db *sql.DB) {
	defer wgSaveComments.Done()
	stmt, err := db.Prepare("INSERT INTO comments(ID, post_id, name, email, body) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(comment.Id, comment.PostId, comment.Name, comment.Email, comment.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func main() {
	db, err := sql.Open("mysql",
		"root:qwerty@tcp(127.0.0.1:3306)/testdb")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	getUserPosts(7)

	for _, post := range postsList {
		isSuccess := savePostToDB(post, db)
		if isSuccess {
			wgGetComments.Add(1)
			go getPostComments(post.Id, db)
		}
	}

	wgGetComments.Wait()
	wgSaveComments.Wait()

}
