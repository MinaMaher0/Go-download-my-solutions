package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Enter Codeforces handle : ")
	var handle string
	fmt.Scanf("%s", &handle)

	createSolutionDir()

	problemSolutionsLinks := fetchSolutionsLink(handle)

	downloadSolutions(problemSolutionsLinks)
}

func createSolutionDir() {
	currentPath, _ := os.Getwd()
	os.Mkdir(currentPath+"/solutions/", 0777)
}

func getCode(solutionLink string) string {
	c := colly.NewCollector()
	var code = ""
	c.OnHTML("pre", func(e *colly.HTMLElement) {
		code += e.Text

	})
	c.Visit(solutionLink)
	return strings.TrimSuffix(code, "?????")
}

func downloadSolutions(problemSolutionsLinks []problemSolution) {

	currentPath, _ := os.Getwd()
	problemSolutionCnt := make(map[string]int)
	for idx, prolbem := range problemSolutionsLinks {
		cnt, state := problemSolutionCnt[prolbem.name]
		if state {
			f, _ := os.Create(currentPath + "/solutions/" + prolbem.name + "-" + strconv.Itoa(cnt))
			f.WriteString(getCode(prolbem.solutionLink))
		} else {
			f, _ := os.Create(currentPath + "/solutions/" + prolbem.name)
			f.WriteString(getCode(prolbem.solutionLink))
		}
		fmt.Println("Loading...", strconv.Itoa(idx+1)+"/"+strconv.Itoa(len(problemSolutionsLinks)))
		problemSolutionCnt[prolbem.name]++
	}
}

func fetchSolutionsLink(handle string) []problemSolution {
	response, err := http.Get("https://codeforces.com/api/user.status?handle=" + handle + "&from=1&count=100000")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	solutionsLinks := make([]problemSolution, 0)
	if responseObject.Status == "OK" {
		for _, submission := range responseObject.Submissions {
			if submission.Verdict == "OK" && submission.ContestID < 100000 {
				solutionLink := "https://codeforces.com/contest/" + strconv.Itoa(submission.ContestID) + "/submission/" + strconv.Itoa(submission.ID)
				solutionsLinks = append(solutionsLinks, problemSolution{
					name:         submission.Problem.Index + " - " + submission.Problem.Name,
					solutionLink: solutionLink,
				})
			}
		}
	} else {
		fmt.Println("Handle not found")
		os.Exit(1)
	}

	return solutionsLinks
}

type problemSolution struct {
	name         string
	solutionLink string
}

type Response struct {
	Status      string       `json:"status"`
	Submissions []Submission `json:"result"`
}

type Submission struct {
	ID                  int     `json:"id"`
	ContestID           int     `json:"contestId"`
	Problem             Problem `json:"problem"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`
}

type Problem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
}
