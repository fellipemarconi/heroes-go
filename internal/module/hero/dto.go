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
