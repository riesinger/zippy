package main

import (
	"encoding/json"
	"fmt"
	"github.com/arial7/zippy/database"
	"github.com/aymerick/raymond"
	"github.com/gosimple/slug"
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"net/http"
	"os"
)

const VERSION = "0.0.1"

var logger zap.Logger
var db *database.MgoAdapter

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
	// TODO: Remove this at least in 1.0.0, as it is only used to enable debugging
	viper.SetDefault("environment", "development")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/viper")
	viper.AddConfigPath(".")

	// Load the actual config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error while reading config: %s\n", err))
	}

	// Setup logging
	if viper.GetString("environment") == "development" {
		logger = zap.New(zap.NewTextEncoder(), zap.Output(os.Stdout))
		logger.SetLevel(zap.DebugLevel)
	} else {
		logfile, err := os.OpenFile(viper.GetString("server.logDir")+"/log",
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(fmt.Errorf("Failed to create logfile: %s\n", err))
		}
		defer logfile.Close()
		logger = zap.New(zap.NewTextEncoder(), zap.Output(logfile))
		logger.SetLevel(zap.InfoLevel)
	}

	themeDir := "themes/" + viper.GetString("site.theme")

	db = database.NewMgoAdapter(logger)
	err = db.Dial(viper.GetString("database.url"), viper.GetString("database.name"))
	if err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	checkThemeExists(themeDir)
	loadTemplates(themeDir)

	setupRoutes(themeDir)
}

func checkThemeExists(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Fatal("Requested theme does not exist", zap.String("directory", path))
		os.Exit(1)
	} else if err != nil {
		logger.Error("Could not check for theme", zap.String("directory", path), zap.Error(err))
	}
}
