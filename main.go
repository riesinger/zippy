package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/gosimple/slug"
	"github.com/russross/blackfriday"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	ErrorText string      `json:"errorText"`
	Payload   interface{} `json:"payload"`
}

type SaveArticleRequest struct {
	Title        string `json:"title"`
	IsDraft      bool   `json:"isDraft"`
	MarkdownBody string `json:"body"`
}

type Article struct {
	Slug         string
	Title        string
	MarkdownBody string
	Body         string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	// Decode the request
	decoder := json.NewDecoder(r.Body)
	var articleReq SaveArticleRequest
	err := decoder.Decode(&articleReq)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	// Early response to client
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{true, "", nil}
	json.NewEncoder(w).Encode(response)

	// Compile the article content
	pageSlug := slug.Make(articleReq.Title)
	htmlContent := string(blackfriday.MarkdownCommon([]byte(articleReq.MarkdownBody)))

	article := Article{
		Slug:         pageSlug,
		Title:        articleReq.Title,
		MarkdownBody: articleReq.MarkdownBody,
		Body:         htmlContent,
	}

	articles = append(articles, article)

}

func renderPost(w http.ResponseWriter, r *http.Request) {
	a := articles[0]
	ctx := map[string]string{
		"title":   a.Title,
		"content": a.Body,
	}
	result, err := postTemplate.Exec(ctx)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(result))

}

var articles []Article
var postTemplate *raymond.Template

func main() {
	// Load the "post" handlebar remplate and compile it
	tmpl, err := ioutil.ReadFile("static/templates/post.hbs")
	if err != nil {
		panic(err)
	}

	postTemplate, err = raymond.Parse(string(tmpl))
	if err != nil {
		panic(err)
	}

	// API routes
	http.HandleFunc("/api/saveArticle", saveArticle)

	// Actual blog routes
	http.HandleFunc("/posts", renderPost)

	// Static files
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/js", http.FileServer(http.Dir("static/js")))
	http.Handle("/css", http.FileServer(http.Dir("static/css")))
	http.ListenAndServe(":8080", nil)
}
