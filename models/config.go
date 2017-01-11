package models

type Configuration struct {
	SiteName  string `bson:"siteName" json:"siteName"`
	ThemeName string `bson:"themeName" json:"themeName"`
	BaseURL   string `bson:"baseURL" json:"baseURL"`
	IsSetup   bool   `bson:"isSetup" json:"isSetup"`
}
