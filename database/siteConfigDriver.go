package database

import (
	"github.com/arial7/zippy/models"
)

type SiteConfigDriver interface {
	GetSiteConfig() *models.Configuration
	SetInitialSiteConfig(*models.SignupData) error
}
