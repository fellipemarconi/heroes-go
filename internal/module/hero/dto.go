package hero

import "mime/multipart"

type CreateHeroInput struct {
	Name        string   `validate:"required,min=2,max=100"`
	Slug        string   `validate:"required,min=2,max=30"`
	Universe    string   `validate:"required,min=2,max=50"`
	Powers      []string `validate:"required,min=1,dive,min=2,max=50"`
	Alignment   string   `validate:"required,oneof=hero villain anti-hero"`
	Description string   `validate:"max=255"`
}

type ListHeroesInput struct {
	Universe   string
	Alignment  string
	IsActive   *bool
	PageOffset int32
	PageSize   int32
}

type Hero struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Alignment   string   `json:"alignment"`
	Universe    string   `json:"universe"`
	Powers      []string `json:"powers"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
	Image       string   `json:"image"`
}

type UpdateHeroStatusInput struct {
	HeroID   string `validate:"required,uuid4"`
	IsActive *bool  `validate:"required, boolean"`
}

type UpdateHeroImageInput struct {
	File        multipart.File
	FileHeader  *multipart.FileHeader
	ContentType string
}

type SearchHeroesInput struct {
	Query      string `form:"q" validate:"required,min=2,max=100"`
	Universe   string `form:"universe"`
	Alignment  string `form:"alignment"`
	PageOffset int32  `form:"offset" validate:"min=0"`
	PageSize   int32  `form:"limit" validate:"min=1,max=100"`
}

type SearchHeroResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Universe    string   `json:"universe"`
	Alignment   string   `json:"alignment"`
	Powers      []string `json:"powers"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Rank        float32  `json:"rank"`
}

type SearchResultsResponse struct {
	Results []SearchHeroResponse `json:"results"`
	Total   int64                `json:"total"`
	Limit   int32                `json:"limit"`
	Offset  int32                `json:"offset"`
}
