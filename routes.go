package main

import (
	"github.com/aymerick/raymond"
	"github.com/spf13/viper"
	"github.com/uber-go/zap"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	tmpl, err := getTemplate("post", false)
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

func renderPage(tmpl *raymond.Template, context interface{}) ([]byte, error) {
	// We got the template, now compile the article into it!
	result, err := tmpl.Exec(context)
	return []byte(result), err
}

// setupRoutes() binds the handlers to specific paths and starts the server
func setupRoutes(themeDir string) {

	staticPath, _ := filepath.Abs(filepath.Join(themeDir, "static"))
	h := http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath)))

	http.HandleFunc("/", standardPageHandler)
	http.HandleFunc("/adm", adminPageHandler)
	http.Handle("/static/", h)

	// Start the server
	logger.Info("Server started", zap.String("port", viper.GetString("server.port")))
	http.ListenAndServe(":"+viper.GetString("server.port"), nil)
}
