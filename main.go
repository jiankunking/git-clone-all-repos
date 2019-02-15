package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

type repo struct {
	RepoName string `json:"name"`
}

func main() {
	//根据机器情况自己调节
	runtime.GOMAXPROCS(2)
	var (
		repos       []repo
		username    string
		totalNumber int
	)
	fmt.Println("Enter github username:")
	fmt.Scanf("%s", &username)
	fmt.Println("Enter repos total number:")
	fmt.Scanf("%d", &totalNumber)
	//username = "rancher"
	//totalNumber = 1000

	cycles := math.Ceil(float64(totalNumber / 100))
	fmt.Println("cycles", cycles)
	for cycle := 1; cycle <= int(cycles); cycle++ {
		url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=100", username, cycle)
		fmt.Println("Fetching data from github ... ")
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			break
		}
		var b bytes.Buffer
		io.Copy(&b, resp.Body)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		var tmp []repo
		json.Unmarshal([]byte(b.String()), &tmp)
		if len(tmp) == 0 {
			break
		}
		repos = append(repos, tmp...)
	}

	fmt.Println("Numbers of repositories:", len(repos))
	var wg sync.WaitGroup

	for i := 0; i < len(repos); i++ {
		repoName := repos[i].RepoName
		wg.Add(1)
		go func(username, repoName string) {
			defer wg.Done()
			RepoURL := fmt.Sprintf("https://github.com/%s/%s.git", username, repoName)
			FilePATH := fmt.Sprintf("%s/%s", username, repoName)
			fmt.Printf("Cloning %d %s repository... \n", i, repoName)
			cmd := exec.Command("git", "clone", RepoURL, FilePATH)
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Finished Clone %s repository\n", repoName)
		}(username, repoName)
	}
	wg.Wait()
	fmt.Println("Finished Clone all repository")
}
