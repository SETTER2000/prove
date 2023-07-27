package entity

type (
	Balance struct {
		UserID  string  `json:"-"`
		Current float32 `json:"current"`
	}
)
