package models

type SignupData struct {
	SiteName        string `json:"siteName"`
	BaseURL         string `json:"baseURL"`
	Theme           string `json:"theme"`
	AccountName     string `json:"accountName"`
	AccountMail     string `json:"accountMail"`
	AccountPassword string `json:"accountPassword"`
}
