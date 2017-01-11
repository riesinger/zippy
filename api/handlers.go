package api

import (
	"encoding/json"
	"github.com/arial7/zippy/database"
	"github.com/arial7/zippy/models"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/uber-go/zap"
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

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	collection, hasCollection := vars["collection"]
	logger.Debug("Requested article", zap.String("action", action), zap.String("collection", collection))

	// TODO:
	_ = hasCollection

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

		w.Write(SuccessResponse().Marshal())

		db.CreateArticle(a)
	case "get":
		if hasCollection == false {
			w.Write(ErrorResponse("You need to specify an article to get", ErrNoCollection).Marshal())
			return
		}
		articles, err := db.GetArticleBySlug(collection)
		if err != nil {
			logger.Error("Could not get article for slug", zap.Error(err))
			w.WriteHeader(500)
			return
		}
		w.Write(PayloadResponse(articles).Marshal())
	}
}
