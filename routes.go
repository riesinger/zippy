package main

import (
	"github.com/arial7/zippy/api"
	"github.com/aymerick/raymond"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func standardPageHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Requested standard page", zap.Object("url", r.URL.Path))

	a, err := db.GetArticleByFullPath(r.URL.Path)
	if err == mgo.ErrNotFound {
		// TODO: Show a proper 404 page
		w.WriteHeader(404)
		return
	}

	context := map[string]string{
		"title":   a.Title,
		"content": a.HtmlBody,
	}
	tmpl, err := getTemplate(a.Template)
	if err != nil {
		logger.Error("Could not get template", zap.String("template", a.Template), zap.Error(err))
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
	f, _ := ioutil.ReadFile(filepath.Join("themes", viper.GetString("site.theme"),
		"admin", "editor.html"))
	w.Write(f)
}

func renderPage(tmpl *raymond.Template, context interface{}) ([]byte, error) {
	// We got the template, now compile the article into it!
	result, err := tmpl.Exec(context)
	return []byte(result), err
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
