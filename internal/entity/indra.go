package entity

import (
	"github.com/SETTER2000/prove/config"
)

type (
	Indras []Indra
	UserID string

	// Indra -.
	Indra struct {
		*config.Config `json:"-"`
		Slug           `json:"slug,omitempty" example:"1674872720465761244B_5"`
		URL            `json:"url,omitempty" example:"https://example.com/go/to/home.html"`
		UserID         `json:"user_id,omitempty"`
		Del            bool `json:"del"`
	}

	User struct {
		UserID  string `json:"user_id" example:"1674872720465761244B_5"`
		DelLink []Slug
		Urls    []List
	}
)
