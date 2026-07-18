package dto

type AddSiteRequest struct {
	URL  string `json:"url" example:"https://example.com"`
	Name string `json:"name" example:"Example"`
}
