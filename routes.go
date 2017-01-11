package main

import (
	"encoding/json"
	"github.com/arial7/zippy/api"
	"github.com/arial7/zippy/models"
	"github.com/aymerick/raymond"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

func standardPageHandler(w http.ResponseWriter, r *http.Request) {
	// First, get the page path
	//vars := mux.Vars(r)
	//pagePath := vars["pagePath"]

	//logger.Debug("Requested standard page", zap.String("path", pagePath))

	logger.Debug("Requested standard page", zap.Object("url", r.URL))

	context := map[string]string{
		"title":   "Somepost",
		"content": "Here is some content",
	}
	tmpl, err := getTemplate("post")
	if err != nil {
		logger.Error("Could not get template", zap.Error(err))
	} else {
		response, err := renderPage(tmpl, context)
		if err == nil {
			w.Write(response)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: We should really use the themeDir here!
	// TODO: The editor should also be a HBS template with some variables
	f, _ := ioutil.ReadFile(filepath.Join("themes", viper.GetString("site.theme"), "admin", "editor.html"))
	w.Write(f)
}

// - API HANDLERS -

func newArticleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Requested api page", zap.Object("url", r.URL))
	decoder := json.NewDecoder(r.Body)
	var a models.Article
	err := decoder.Decode(&a)
	if err != nil {
		logger.Error("Cannot decode new article", zap.Error(err))
	}
	defer r.Body.Close()

	db.CreateArticle(&a)
}

func getArticleBySlugHandler(w http.ResponseWriter, r *http.Request) {

}

func getArticlesByPathHandler(w http.ResponseWriter, r *http.Request) {

}

func renderPage(tmpl *raymond.Template, context interface{}) ([]byte, error) {
	// We got the template, now compile the article into it!
	result, err := tmpl.Exec(context)
	return []byte(result), err
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Requested api page", zap.Object("url", r.URL))
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Debug("Requested article", zap.String("action", vars["action"]),
		zap.String("list", vars["list"]))

	switch vars["action"] {
	case "new":
		decoder := json.NewDecoder(r.Body)
		var a models.Article
		err := decoder.Decode(&a)
		if err != nil {
			logger.Error("Cannot decode new article", zap.Error(err))
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()
		a.CreatedAt = time.Now().UTC()
		a.UpdatedAt = time.Now().UTC()

		db.CreateArticle(&a)
	}
}

// setupRoutes() binds the handlers to specific paths and starts the server
func setupRoutes(themeDir string) {
	api.SetupHandlers(db, logger)

	apiRouter := mux.NewRouter()
	staticPath, _ := filepath.Abs(filepath.Join(themeDir, "static"))
	h := http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath)))

	apiRouter.PathPrefix("/article").Path("/{action}").HandlerFunc(api.ArticleHandler)
	apiRouter.PathPrefix("/article").Path("/{action}/{collection}").HandlerFunc(api.ArticleHandler)

	http.HandleFunc("/", standardPageHandler)
	http.HandleFunc("/adm", adminPageHandler)
	http.Handle("/api/", http.StripPrefix("/api", apiRouter))
	http.Handle("/static/", h)

	// Start the server
	logger.Info("Server started", zap.String("port", viper.GetString("server.port")))
	http.ListenAndServe(":"+viper.GetString("server.port"), nil)
}
