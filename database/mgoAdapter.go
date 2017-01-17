package database

import (
	"errors"
	"fmt"
	"github.com/arial7/zippy/models"
	"github.com/arial7/zippy/utils"
	"github.com/uber-go/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const ArticleCollectionName = "articles"
const SiteConfigCollectionName = "siteconfig"
const UserCollectionName = "users"

type MgoAdapter struct {
	ConfigCollection *mgo.Collection
	SiteConfig       *models.Configuration
	Users            *mgo.Collection
	Articles         *mgo.Collection
	Session          *mgo.Session
	Logger           zap.Logger
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
		m.Articles = session.DB(databaseName).C(ArticleCollectionName)
		m.ConfigCollection = session.DB(databaseName).C(SiteConfigCollectionName)
		m.Users = session.DB(databaseName).C(UserCollectionName)
		return nil
	}
}

func (m *MgoAdapter) CreateArticle(article *models.Article) error {
	m.Logger.Debug("Creating new article", zap.String("title", article.Title))
	err := m.Articles.Insert(article)
	if err != nil {
		m.Logger.Error("Could not create article", zap.Error(err))
	}
	return err
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

func (m *MgoAdapter) GetSiteConfig() *models.Configuration {
	if m.SiteConfig == nil {
		err := m.ConfigCollection.Find(nil).One(&m.SiteConfig)
		if err == mgo.ErrNotFound {
			m.Logger.Info("There is no site config, please visit /setup")
			m.SiteConfig = &models.Configuration{IsSetup: false}
		} else if err != nil {
			m.Logger.Fatal("Could not retrive site config", zap.Error(err))
			os.Exit(1)
		}
	}
	return m.SiteConfig
}

func (m *MgoAdapter) SetInitialSiteConfig(data *models.SignupData) error {
	m.Logger.Debug("Setting initial site config")

	if m.SiteConfig.IsSetup == true {
		m.Logger.Error("Tried to set initial config data when already configured")
		return errors.New("set initial config when already configured")
	}

	newConfig := &models.Configuration{
		SiteName:  data.SiteName,
		ThemeName: data.Theme,
		BaseURL:   data.BaseURL,
		IsSetup:   true,
	}

	err := m.ConfigCollection.Insert(newConfig)
	if err != nil {
		m.Logger.Error("Could not update site config", zap.Error(err))
		return err
	}

	m.SiteConfig = newConfig

	return nil
}

func (m *MgoAdapter) CreateOwnerUser(data *models.SignupData) error {
	m.Logger.Debug("Creating owner user")
	// First, look for an existing owner user
	var dummy models.User
	err := m.Users.Find(bson.M{"role": 0}).One(&dummy)
	if err != mgo.ErrNotFound {
		m.Logger.Error("Tried to create owner user", zap.Error(err))
		return errors.New("creating a second owner user is illegal")
	}

	// There is no owner user, so go on with encrypting the password
	user := models.User{
		FullName: data.AccountName,
		Email:    data.AccountMail,
		UID:      utils.NewUUID(),
		Role:     models.UserRoleOwner,
	}

	hashedPassword, salt, err := utils.HashWithNewSalt(data.AccountPassword)

	if err != nil {
		m.Logger.Error("Could not hash password", zap.Error(err))
		return errors.New("failed to hash password")
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = m.Users.Insert(&user)

	if err != nil {
		m.Logger.Error("Could not save owner user", zap.Error(err))
		return errors.New("failed to save owner")
	}

	m.Logger.Debug("Created owner")

	return nil

}

func (m *MgoAdapter) Close() {
	m.Session.Close()
}
