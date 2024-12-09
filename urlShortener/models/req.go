package models

type Req struct {
	LongURL string `json:"long_url" binding:"required,url"`
}
