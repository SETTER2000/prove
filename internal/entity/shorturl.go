package entity

type (

	// Slug -.
	Slug string

	// URL -.
	URL string

	// List -.
	List struct {
		ShortURL URL                                                                 `json:"short_url" example:"1674872720465761244B_5"` // Строковый идентификатор
		URL      `json:"original_url" example:"https://example.com/go/to/home.html"` // URL для сокращения
	}

	// CountUsers кол-во пользователей в сервисе
	CountUsers int

	// CountURLs кол-во сокращённых URL в сервисе
	CountURLs int
)
