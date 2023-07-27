package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	SolutionList []Solution

	Solution struct {
		SolutionID  string `json:"solution_id,omitempty"`
		Description string `json:"description" validate:"required"`
		TaskID      string `json:"task_id,omitempty"`
		Solution    string `json:"solution"  validate:"required"`
		UserID      string `json:"user_id,omitempty"`
		UploadedAt  string `json:"created,omitempty"`
	}

	SolutionData struct {
		UserID   `json:"-"`
		TaskID   string `json:"task_id" validate:"required"`
		Data     []int  `json:"data" validate:"required"`
		Solution []int  `json:"result"`
	}
)

func (c *Solution) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Solution, validation.Required),
	)
}

// Sanitize очищает поля, для того чтоб они не возвращались в ответе
func (c *Solution) Sanitize() {
	c.UserID = ""
}
