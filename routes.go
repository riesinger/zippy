package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/arial7/zippy/api"
	"github.com/aymerick/raymond"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"os"
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
	if isSetup == false {
		http.Redirect(w, r, "/setup", 303)
	}
	// TODO: We should really use the themeDir here!
	// TODO: The editor should also be a HBS template with some variables
	f, _ := ioutil.ReadFile(filepath.Join("themes", db.GetSiteConfig().ThemeName,
		"admin", "login.html"))
	w.Write(f)
}

func setupPageHandler(w http.ResponseWriter, r *http.Request) {
	if isSetup == true {
		http.Redirect(w, r, "/", 303)
	}
	// TODO: Be more intelligent about assuming the default theme for setup
	f, _ := ioutil.ReadFile(filepath.Join("themes", "default",
		"admin", "setup.html"))
	w.Write(f)

}

func renderPage(tmpl *raymond.Template, context interface{}) ([]byte, error) {
	// We got the template, now compile the article into it!
	result, err := tmpl.Exec(context)
	return []byte(result), err
}

// setupRoutes() binds the handlers to specific paths and starts the server
func setupRoutes(themeDir string) {

	// This is used for the initial setup
	if themeDir == "" {
		themeDir = "themes/default"
	}

	api.SetupHandlers(db, logger)

	apiApp := rest.NewApi()

	if viper.GetString("environment") == "development" {
		apiApp.Use(rest.DefaultDevStack...)
	} else {
		apiApp.Use(rest.DefaultProdStack...)
	}

	apiRouter, err := rest.MakeRouter(
		rest.Get("/article/#slug", api.GetArticleHandler),
		rest.Post("/article/new", api.CreateArticleHandler),
		rest.Post("/initialSetup", api.SetupHandler),
	)

	if err != nil {
		logger.Fatal("Could not setup API router", zap.Error(err))
		os.Exit(1)
	}

	apiApp.SetApp(apiRouter)

	staticPath, _ := filepath.Abs(filepath.Join(themeDir, "static"))
	h := http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath)))

	http.HandleFunc("/", standardPageHandler)
	http.HandleFunc("/adm/", adminPageHandler)
	http.HandleFunc("/setup", setupPageHandler)
	http.Handle("/api/", http.StripPrefix("/api", apiApp.MakeHandler()))
	http.Handle("/static/", h)

	// Start the server
	logger.Info("Server started", zap.String("port", viper.GetString("server.port")))
	http.ListenAndServe(":"+viper.GetString("server.port"), nil)
}
