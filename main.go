package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var posts []string
var wg sync.WaitGroup
var mut sync.Mutex

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

func main() {
	start := time.Now()

	getPosts()
	wg.Wait()

	programTime := time.Since(start)
	fmt.Println(posts)
	fmt.Println(programTime)
}
