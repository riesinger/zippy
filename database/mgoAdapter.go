package database

import (
	"errors"
	"fmt"
	"github.com/arial7/zippy/models"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ArticleCollectionName = "articles"

type MgoAdapter struct {
	Articles *mgo.Collection
	Session  *mgo.Session
	Logger   zap.Logger
}

func NewMgoAdapter(logger zap.Logger) *MgoAdapter {
	return &MgoAdapter{
		Logger: logger.With(zap.String("component", "MgoAdapter")),
	}
}

func (m *MgoAdapter) Dial(url string, databaseName string) error {
	m.Logger.Debug("Dialing", zap.String("url", url))
	session, err := mgo.Dial(url)
	if err != nil {
		return errors.New(fmt.Sprintf("error while dialing: %s", err))
	} else {
		m.Session = session
		c := session.DB(databaseName).C(ArticleCollectionName)
		m.Articles = c
		return nil
	}
}

func (m *MgoAdapter) CreateArticle(article *models.Article) {
	m.Logger.Debug("Creating new article", zap.String("title", article.Title))
	err := m.Articles.Insert(article)
	if err != nil {
		m.Logger.Error("Could not create article", zap.Error(err))
	}
}

func (m *MgoAdapter) GetArticleBySlug(slug string) (*models.Article, error) {
	var a *models.Article
	err := m.Articles.Find(bson.M{"slug": slug}).One(&a)
	if err != nil {
		m.Logger.Error("Could not get article by slug", zap.String("slug", slug), zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (m *MgoAdapter) GetArticlesByPath(path string) ([]*models.Article, error) {
	var as []*models.Article
	err := m.Articles.Find(bson.M{"path": path}).All(&as)
	if err != nil {
		m.Logger.Error("Could not get all articles for path", zap.String("path", path), zap.Error(err))
		return nil, err
	}
	return as, nil
}

func (m *MgoAdapter) GetArticleByFullPath(fullPath string) (*models.Article, error) {
	var a *models.Article
	err := m.Articles.Find(bson.M{"fPath": fullPath}).One(&a)
	if err != nil {
		m.Logger.Warn("Could not get article for", zap.String("path", fullPath), zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (m *MgoAdapter) Close() {
	m.Session.Close()
}
