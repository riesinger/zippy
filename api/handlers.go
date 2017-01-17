package api

import (
	"encoding/json"
	"github.com/arial7/zippy/database"
	"github.com/arial7/zippy/models"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/russross/blackfriday"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"time"
)

var logger zap.Logger
var db *database.MgoAdapter

func SetupHandlers(database *database.MgoAdapter, loggerP zap.Logger) {
	db = database
	logger = loggerP.With(zap.String("component", "API"))
}

func unmarshalArticle(body io.ReadCloser) (*models.Article, error) {
	var a models.Article
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&a)
	defer body.Close()
	if err != nil {
		logger.Error("Cannot decode article", zap.Error(err))
		return nil, err
	}
	return &a, nil

}

func SetupHandler(w http.ResponseWriter, r *http.Request) {
	if db.GetSiteConfig().IsSetup {
		w.WriteHeader(400)
		return
	}

	var signupData models.SignupData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&signupData)
	defer r.Body.Close()
	if err != nil {
		logger.Error("Cannot decode signup data", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	err = db.SetInitialSiteConfig(&signupData)
	if err != nil {
		logger.Error("Cannot setup site config data", zap.Error(err))
		w.Write(ErrorResponse("Internal server error", ErrInitialConfig).Marshal())
		w.WriteHeader(500)
		return
	}

	err = db.CreateOwnerUser(&signupData)
	if err != nil {
		logger.Error("Cannot create admin user", zap.Error(err))
		w.Write(ErrorResponse("Internal server error", ErrCreateOwner).Marshal())
		w.WriteHeader(500)
		return
	}

	w.Write(SuccessResponse().Marshal())

}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	collection, hasCollection := vars["collection"]
	logger.Debug("Requested article", zap.String("action", action), zap.String("collection", collection))

	switch action {
	case "new":
		a, err := unmarshalArticle(r.Body)
		if err != nil {
			w.Write(ErrorResponse("Article was invalid", ErrUnmarshalArticle).Marshal())
			return
		}
		a.UpdatedAt = time.Now().UTC()
		a.CreatedAt = time.Now().UTC()
		a.Slug = slug.Make(a.Title)
		a.FullPath = a.Path + a.Slug

		// TODO: Sanetize the input with bluemonday
		a.HtmlBody = string(blackfriday.MarkdownCommon([]byte(a.MarkdownBody)))

		// TODO: We should check if the article *can* be created and then send a success response
		// immediately (clients don't care for backend processing, so do this silently)
		w.Write(SuccessResponse().Marshal())
		db.CreateArticle(a)

	case "get":
		if hasCollection == false {
			w.Write(ErrorResponse("You need to specify an article to get", ErrNoCollection).Marshal())
			return
		}
		articles, err := db.GetArticleBySlug(collection)
		if err == mgo.ErrNotFound {
			w.Write(ErrorResponse("Article not found", ErrArticleNotFound).Marshal())
			return
		}
		if err != nil {
			logger.Error("Could not get article for slug", zap.Error(err))
			w.Write(ErrorResponse("Could not get article", ErrInternal).Marshal())
			w.WriteHeader(500)
			return
		}
		w.Write(PayloadResponse(articles).Marshal())
	case "update":
		if hasCollection == false {
			w.Write(ErrorResponse("You need to specify an article to update", ErrNoCollection).Marshal())
			return
		}
		logger.Warn("Updating articles is not implemented yet")
	}
}
