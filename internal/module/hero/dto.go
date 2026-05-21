package hero

type CreateHeroInput struct {
	Name        string   `validate:"required,min=2,max=100"`
	Slug        string   `validate:"required,min=2,max=30"`
	Universe    string   `validate:"required,min=2,max=50"`
	Powers      []string `validate:"required,min=1,dive,min=2,max=50"`
	Alignment   string   `validate:"required,oneof=hero villain anti-hero"`
	Description string   `validate:"max=255"`
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
}
