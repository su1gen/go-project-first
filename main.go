package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getPosts() {
	response, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Println("wrong request")
		return
	}

	posts, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error with response")
		return
	}
	//Convert the body to type string
	postsForPrint := string(posts)
	fmt.Println(postsForPrint)
}

func main() {
	getPosts()
}
