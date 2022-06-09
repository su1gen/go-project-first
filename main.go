package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var posts []string
var wg sync.WaitGroup
var mut sync.Mutex

func createPostFiles(postsSlice []string, path string) {
	for index, post := range postsSlice {
		newFile, err := os.Create(path + strconv.Itoa(index+1) + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		//newFile.WriteString(post)
		ioutil.WriteFile(path+strconv.Itoa(index+1)+".txt", []byte(post), 0644)
		newFile.Close()

	}
}

func getOnePost(index int) {
	defer wg.Done()
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", index+1)
	response, _ := http.Get(url)
	post, _ := ioutil.ReadAll(response.Body)
	mut.Lock()
	defer mut.Unlock()
	posts = append(posts, string(post))
}

func getPosts() {
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go getOnePost(i)
	}
}

func removePostFiles(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		os.Remove(path + file.Name())
	}
}

func main() {
	start := time.Now()
	path := "./storage/posts/"

	removePostFiles(path)

	getPosts()
	wg.Wait()

	createPostFiles(posts[:5], path)

	programTime := time.Since(start)
	fmt.Println(programTime)
}
