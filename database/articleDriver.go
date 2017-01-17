package database

import (
	"github.com/arial7/zippy/models"
)

type ArticleDriver interface {
	CreateArticle(article *models.Article) error
	GetArticleBySlug(slug string) (*models.Article, error)
	GetArticlesByPath(path string) ([]*models.Article, error)
	GetArticleByFullPath(fullPath string) (*models.Article, error)
}
