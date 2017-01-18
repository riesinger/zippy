package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/arial7/zippy/database"
	"github.com/arial7/zippy/models"
	"github.com/gosimple/slug"
	"github.com/russross/blackfriday"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"time"
)

var logger zap.Logger
var db *database.MgoAdapter

func SetupHandlers(database *database.MgoAdapter, loggerP zap.Logger) {
	db = database
	logger = loggerP.With(zap.String("component", "API"))
}

func SetupHandler(w rest.ResponseWriter, r *rest.Request) {
	if db.GetSiteConfig().IsSetup {
		w.WriteHeader(400)
		return
	}

	signupData := models.SignupData{}
	err := r.DecodeJsonPayload(&signupData)
	if err != nil {
		logger.Error("Cannot decode setup data", zap.Error(err))
		w.WriteJson(ErrorResponse("Cannot decode setup data", ErrSignupData))
		w.WriteHeader(500)
		return
	}

	err = db.SetInitialSiteConfig(&signupData)
	if err != nil {
		logger.Error("Cannot setup site config data", zap.Error(err))
		w.WriteJson(ErrorResponse("Internal server error", ErrInitialConfig))
		w.WriteHeader(500)
		return
	}

	err = db.CreateOwnerUser(&signupData)
	if err != nil {
		logger.Error("Cannot create admin user", zap.Error(err))
		w.WriteJson(ErrorResponse("Internal server error", ErrCreateOwner))
		w.WriteHeader(500)
		return
	}

	w.WriteJson(SuccessResponse())
}

func CreateArticleHandler(w rest.ResponseWriter, r *rest.Request) {
	a := models.Article{}
	err := r.DecodeJsonPayload(&a)
	if err != nil {
		w.WriteJson(ErrorResponse("Article was invalid", ErrUnmarshalArticle))
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
	w.WriteJson(SuccessResponse())
	db.CreateArticle(&a)

}

func GetArticleHandler(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	logger.Debug("Requested article", zap.String("slug", slug))

	article, err := db.GetArticleBySlug(slug)
	if err == mgo.ErrNotFound {
		w.WriteJson(ErrorResponse("Article not found", ErrArticleNotFound))
		return
	}
	if err != nil {
		logger.Error("Could not get article for slug", zap.Error(err))
		w.WriteJson(ErrorResponse("Could not get article", ErrInternal))
		w.WriteHeader(500)
		return
	}
	w.WriteJson(PayloadResponse(article))
}
