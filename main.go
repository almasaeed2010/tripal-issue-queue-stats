package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Label ...
type Label struct {
	ID      int    `json:"id"`
	NodeID  string `json:"node_id"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Default bool   `json:"default"`
}

// User ...
// This is the user that gets returned in "user" key of the github response
type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingUTL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// PR ...
type PR struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}

// Issue ...
type Issue struct {
	URL               string        `json:"url"`
	RepositoryURL     string        `json:"repository_url"`
	LabelsURL         string        `json:"labels_url"`
	CommentsURL       string        `json:"comments_url"`
	EventsURL         string        `json:"events_url"`
	HTMLURL           string        `json:"html_url"`
	ID                int           `json:"id"`
	NodeID            string        `json:"node_id"`
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	User              User          `json:"user"`
	Labels            []Label       `json:"labels"`
	State             string        `json:"state"`
	Locked            bool          `json:"locked"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	Milestone         interface{}   `json:"milestone"`
	Comments          int           `json:"comments"`
	CreatedAt         string        `json:"created_at"`
	UpdatedAt         string        `json:"updated_at"`
	ClosedAt          interface{}   `json:"closed_at"`
	AuthorAssociation string        `json:"author_association"`
	PullRequest       PR            `json:"pull_request"`
	Body              string        `json:"body"`
}

type Request struct {
	Username string
	Password string
	URL      string
}

// Parse ...
func Parse(info Request) ([]Issue, error) {
	username := info.Username
	password := info.Password
	url := info.URL

	client := http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)

	resp, e := client.Do(req)

	if e != nil {
		fmt.Println("Error in request")
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)

	if e != nil {
		fmt.Println("Error in reading body")
	}

	if !json.Valid(body) {
		fmt.Println("The response was not valid")
		ioutil.WriteFile("response.json", body, 0644)
	}

	var data []Issue

	err := json.Unmarshal([]byte(body), &data)

	return data, err
}

// The main function
func main() {
	username := flag.String("u", "", "Github Username")
	password := flag.String("p", "", "Github Password")

	flag.Parse()

	page := 1
	data, _ := Parse(Request{Username: *username, Password: *password, URL: "https://api.github.com/repos/tripal/tripal/issues/comments?state=all"})

	counts := make(map[string]int)

	for len(data) > 0 {
		for _, datum := range data {
			if value, ok := counts[datum.User.Login]; ok {
				counts[datum.User.Login] = value + 1
			} else {
				counts[datum.User.Login] = 1
			}
		}

		page = page + 1
		data, _ = Parse(Request{
			Username: *username,
			Password: *password,
			URL:      "https://api.github.com/repos/tripal/tripal/issues/comments?state=all&page=" + strconv.Itoa(page)})
	}

	for name, commentsCount := range counts {
		fmt.Println(name, commentsCount)
	}
}
